package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/mvg-fi-dev/bridge/internal/config"
	"github.com/mvg-fi-dev/bridge/internal/db"
	"github.com/mvg-fi-dev/bridge/internal/mixin"
)

const cursorKey = "mixin.snapshots.offset"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.Open(cfg.SQLitePath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.SQL.Close()
	if err := db.Migrate(dbConn.SQL); err != nil {
		log.Fatal(err)
	}

	ksPath := os.Getenv("MIXIN_KEYSTORE_PATH")
	if ksPath == "" {
		log.Fatal("MIXIN_KEYSTORE_PATH is required for snapshot polling (keystore json from Mixin dashboard)")
	}
	ksBytes, err := os.ReadFile(ksPath)
	if err != nil {
		log.Fatalf("read keystore: %v", err)
	}
	ks, err := mixin.ParseSafeKeystore(ksBytes)
	if err != nil {
		log.Fatalf("parse keystore: %v", err)
	}
	client := mixin.NewSDKClient(ks)
	state := db.NewStateRepo(dbConn.SQL)
	snapRepo := db.NewSnapshotsRepo(dbConn.SQL)
	ordersRepo := db.NewOrdersRepo(dbConn.SQL)

	interval := 3 * time.Second
	if v := os.Getenv("MIXIN_POLL_INTERVAL_MS"); v != "" {
		if ms, err := time.ParseDuration(v + "ms"); err == nil {
			interval = ms
		}
	}

	limit := 200

	log.Printf("bridge-worker polling mixin snapshots every %s", interval)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		offset, _, _ := state.Get(ctx, cursorKey)

		snaps, err := client.ListSafeSnapshots(ctx, limit, offset)
		if err != nil {
			log.Printf("poll error: %v", err)
			cancel()
			time.Sleep(interval)
			continue
		}

		// Mixin snapshots are usually returned in reverse-chronological order.
		// We'll process from oldest to newest within this batch.
		for i := len(snaps) - 1; i >= 0; i-- {
			s := snaps[i]
			internalSnap := mixin.SafeSnapshotToInternal(s)
			raw, _ := json.Marshal(map[string]any{"data": internalSnap})
			inserted, err := snapRepo.InsertIfNew(ctx, s.SnapshotID, time.Now().UTC(), string(raw), internalSnap)
			if err != nil {
				log.Printf("insert snapshot err=%v", err)
				continue
			}
			if inserted {
				// Try match order by memo for Mixin-internal payments.
				if s.Memo != "" {
					creditedAt := s.CreatedAt.UTC()
					_, _ = ordersRepo.SetDepositCreditedByMemo(ctx, s.Memo, s.SnapshotID, creditedAt, s.Amount, s.AssetID)
				}
			}

			// advance cursor to newest snapshot id we saw
			offset = s.SnapshotID
		}

		if offset != "" {
			_ = state.Set(ctx, cursorKey, offset)
		}

		cancel()
		time.Sleep(interval)
	}
}
