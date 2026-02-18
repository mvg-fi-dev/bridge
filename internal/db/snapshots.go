package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/mvg-fi-dev/bridge/internal/mixin"
)

type SnapshotsRepo struct{ DB *sql.DB }

func NewSnapshotsRepo(db *sql.DB) *SnapshotsRepo { return &SnapshotsRepo{DB: db} }

// InsertIfNew stores the raw snapshot JSON (caller should pass raw JSON string).
// Returns true if inserted, false if already existed.
func (r *SnapshotsRepo) InsertIfNew(ctx context.Context, snapshotID string, receivedAt time.Time, rawJSON string, s *mixin.Snapshot) (bool, error) {
	createdAt := sql.NullString{}
	if s != nil {
		if t, err := s.CreatedAtTime(); err == nil && t != nil {
			createdAt = sql.NullString{String: t.UTC().Format(time.RFC3339Nano), Valid: true}
		}
	}

	res, err := r.DB.ExecContext(ctx, `
INSERT OR IGNORE INTO mixin_snapshots(
  snapshot_id, received_at, raw_json,
  created_at, amount, asset_id, opponent_id, memo
) VALUES(?,?,?,?,?,?,?,?)
`, snapshotID, receivedAt.UTC().Format(time.RFC3339Nano), rawJSON,
		createdAt,
		nullStr(sAmount(s)), nullStr(sAssetID(s)), nullStr(sOpponentID(s)), nullStr(sMemo(s)),
	)
	if err != nil {
		return false, err
	}
	a, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return a > 0, nil
}

func nullStr(v string) interface{} {
	if v == "" {
		return nil
	}
	return v
}

func sAmount(s *mixin.Snapshot) string {
	if s == nil { return "" }
	return s.Amount
}
func sAssetID(s *mixin.Snapshot) string {
	if s == nil { return "" }
	return s.AssetID
}
func sOpponentID(s *mixin.Snapshot) string {
	if s == nil { return "" }
	return s.OpponentID
}
func sMemo(s *mixin.Snapshot) string {
	if s == nil { return "" }
	return s.Memo
}
