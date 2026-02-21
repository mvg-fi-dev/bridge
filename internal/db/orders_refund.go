package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/mvg-fi-dev/bridge/internal/models"
)

func (r *OrdersRepo) MarkRefundingWithDetails(ctx context.Context, orderID, refundAssetID, refundAmount, refundReceivedSnapshotID string) error {
	_, err := r.DB.ExecContext(ctx, `
UPDATE orders
SET
  status = ?,
  refund_asset_id = COALESCE(refund_asset_id, ?),
  refund_amount = COALESCE(refund_amount, ?),
  refund_received_snapshot_id = COALESCE(refund_received_snapshot_id, ?),
  updated_at = ?
WHERE id = ? AND status IN (?, ?)
`,
		string(models.StatusRefunding),
		refundAssetID,
		refundAmount,
		refundReceivedSnapshotID,
		time.Now().UTC().Format(time.RFC3339Nano),
		orderID,
		string(models.StatusExecutingSwap),
		string(models.StatusDepositCredited),
	)
	return err
}

func (r *OrdersRepo) ListRefunding(ctx context.Context, limit int) ([]*models.Order, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.DB.QueryContext(ctx, `
SELECT
  id, public_id, status, created_at, updated_at,
  source_chain, source_asset, amount_in, target_chain, target_asset, target_address,
  estimated_out, min_out, quote_expiry_at,
  pay_window_seconds,
  mixin_opponent_id, mixin_asset_id, mixin_pay_memo, mixin_pay_url,
  deposit_txid, deposit_tx_detected_at, deposit_credited_at, amount_credited, refund_to_address,
  final_out, swap_ref, exinswap_trace_id, withdraw_txid, refund_txid,
  refund_asset_id, refund_amount, refund_received_snapshot_id
FROM orders
WHERE status = ? AND (refund_txid IS NULL OR refund_txid = '')
ORDER BY updated_at ASC
LIMIT ?
`, string(models.StatusRefunding), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// scanOrder scans the common orders row layout used by list/get queries.
// rows can be *sql.Row or *sql.Rows via the RowScanner interface.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanOrder(rs rowScanner) (*models.Order, error) {
	var o models.Order
	var status string
	var createdAt, updatedAt string
	var quoteExpiry sql.NullString
	var depositTxID, depositDetectedAt, depositCreditedAt sql.NullString
	var amountCredited, refundToAddress sql.NullString
	var finalOut, swapRef, exinTrace, withdrawTxID, refundTxID sql.NullString
	var refundAssetID, refundAmount, refundReceivedSnapshotID sql.NullString

	if err := rs.Scan(
		&o.ID, &o.PublicID, &status, &createdAt, &updatedAt,
		&o.SourceChain, &o.SourceAsset, &o.AmountIn, &o.TargetChain, &o.TargetAsset, &o.TargetAddress,
		&o.EstimatedOut, &o.MinOut, &quoteExpiry,
		&o.PayWindowSeconds,
		&o.MixinOpponentID, &o.MixinAssetID, &o.MixinPayMemo, &o.MixinPayURL,
		&depositTxID, &depositDetectedAt, &depositCreditedAt, &amountCredited, &refundToAddress,
		&finalOut, &swapRef, &exinTrace, &withdrawTxID, &refundTxID,
		&refundAssetID, &refundAmount, &refundReceivedSnapshotID,
	); err != nil {
		return nil, err
	}

	o.Status = models.OrderStatus(status)
	if t, err := time.Parse(time.RFC3339Nano, createdAt); err == nil {
		o.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339Nano, updatedAt); err == nil {
		o.UpdatedAt = t
	}
	if quoteExpiry.Valid {
		if t, err := time.Parse(time.RFC3339Nano, quoteExpiry.String); err == nil {
			o.QuoteExpiryAt = &t
		}
	}
	if depositTxID.Valid {
		o.DepositTxID = &depositTxID.String
	}
	if depositDetectedAt.Valid {
		if t, err := time.Parse(time.RFC3339Nano, depositDetectedAt.String); err == nil {
			o.DepositTxDetectedAt = &t
		}
	}
	if depositCreditedAt.Valid {
		if t, err := time.Parse(time.RFC3339Nano, depositCreditedAt.String); err == nil {
			o.DepositCreditedAt = &t
		}
	}
	if amountCredited.Valid {
		o.AmountCredited = &amountCredited.String
	}
	if refundToAddress.Valid {
		o.RefundToAddress = &refundToAddress.String
	}
	if finalOut.Valid {
		o.FinalOut = &finalOut.String
	}
	if swapRef.Valid {
		o.SwapRef = &swapRef.String
	}
	if exinTrace.Valid {
		o.ExinSwapTraceID = &exinTrace.String
	}
	if withdrawTxID.Valid {
		o.WithdrawTxID = &withdrawTxID.String
	}
	if refundTxID.Valid {
		o.RefundTxID = &refundTxID.String
	}
	if refundAssetID.Valid {
		o.RefundAssetID = &refundAssetID.String
	}
	if refundAmount.Valid {
		o.RefundAmount = &refundAmount.String
	}
	if refundReceivedSnapshotID.Valid {
		o.RefundReceivedSnapshotID = &refundReceivedSnapshotID.String
	}

	return &o, nil
}
