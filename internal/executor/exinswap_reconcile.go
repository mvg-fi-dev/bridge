package executor

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mvg-fi-dev/bridge/internal/db"
	"github.com/mvg-fi-dev/bridge/internal/exinswap"
	"github.com/mvg-fi-dev/bridge/internal/mixin"
	"github.com/mvg-fi-dev/bridge/internal/models"
)

// ReconcileExinSwapSnapshots scans new snapshots and updates orders:
// - For incoming credits from ExinSwap bot, parse server memo to map TRACE back to our order.
// - Mark refunded / mark withdrawing (WIP: we don't withdraw yet; just persist final_out + status).
//
// NOTE: This is a polling, best-effort reconciler; it should be idempotent.
type ReconcileExinSwapSnapshots struct {
	Orders *db.OrdersRepo
}

func NewReconcileExinSwapSnapshots(orders *db.OrdersRepo) *ReconcileExinSwapSnapshots {
	return &ReconcileExinSwapSnapshots{Orders: orders}
}

func (r *ReconcileExinSwapSnapshots) HandleSnapshot(ctx context.Context, s *mixin.Snapshot) {
	if s == nil {
		return
	}
	// We only care about inbound credits from ExinSwap bot.
	if s.OpponentID != ExinSwapBotUserID {
		return
	}
	if s.Memo == "" {
		return
	}

	memo, err := exinswap.ParseServerMemo(s.Memo)
	if err != nil {
		return
	}

	// memo.Trace should be our transfer trace id. But we used deterministic UUID of order.ID.
	// We need to map trace -> order. For now we store swap_ref as trace id to allow lookup.
	// Better: add a dedicated column exinswap_trace_id.
	o, err := r.findOrderBySwapRef(ctx, memo.Trace)
	if err != nil || o == nil {
		return
	}

	// Update order status based on memo.Type.
	// For swap: success RL means ExinSwap released target asset to us.
	// RF means refund.
	if memo.Type == "RF" {
		log.Printf("exinswap refund order=%s trace=%s", o.PublicID, memo.Trace)
		_ = r.Orders.MarkRefunded(ctx, o.ID, s.SnapshotID)
		return
	}
	if memo.Type == "RL" {
		finalOut := s.Amount
		log.Printf("exinswap release order=%s trace=%s out=%s asset=%s", o.PublicID, memo.Trace, finalOut, s.AssetID)
		// Save final_out and swap_ref; mark withdrawing next.
		_ = r.Orders.MarkWithdrawing(ctx, o.ID, memo.Trace, finalOut)
		return
	}
}

func (r *ReconcileExinSwapSnapshots) findOrderBySwapRef(ctx context.Context, swapRef string) (*models.Order, error) {
	row := r.Orders.DB.QueryRowContext(ctx, `
SELECT id FROM orders WHERE swap_ref = ? LIMIT 1
`, swapRef)
	var id string
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.Orders.GetByID(ctx, id)
}

// Helper to set swap_ref after submitting transfer; this enables reconciliation.
func SetSwapRefOnOrder(ctx context.Context, orders *db.OrdersRepo, orderID string, swapRef string) error {
	_, err := orders.DB.ExecContext(ctx, `
UPDATE orders SET swap_ref = ?, updated_at = ? WHERE id = ?
`, swapRef, time.Now().UTC().Format(time.RFC3339Nano), orderID)
	if err != nil {
		return fmt.Errorf("set swap_ref: %w", err)
	}
	return nil
}
