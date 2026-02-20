package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/mvg-fi-dev/bridge/internal/models"
)

type OrdersRepo struct {
	DB *sql.DB
}

func NewOrdersRepo(db *sql.DB) *OrdersRepo { return &OrdersRepo{DB: db} }

func (r *OrdersRepo) Insert(ctx context.Context, o *models.Order) error {
	_, err := r.DB.ExecContext(ctx, `
INSERT INTO orders (
  id, public_id, status, created_at, updated_at,
  source_chain, source_asset, amount_in, target_chain, target_asset, target_address,
  estimated_out, min_out, quote_expiry_at,
  pay_window_seconds,
  mixin_opponent_id, mixin_asset_id, mixin_pay_memo, mixin_pay_url
) VALUES (?,?,?,?,?, ?,?,?,?,?,?, ?,?,?, ?,?,?,?,?)
`,
		o.ID, o.PublicID, string(o.Status), o.CreatedAt.Format(time.RFC3339Nano), o.UpdatedAt.Format(time.RFC3339Nano),
		o.SourceChain, o.SourceAsset, o.AmountIn, o.TargetChain, o.TargetAsset, o.TargetAddress,
		o.EstimatedOut, o.MinOut, nullableTime(o.QuoteExpiryAt),
		o.PayWindowSeconds,
		o.MixinOpponentID, o.MixinAssetID, o.MixinPayMemo, o.MixinPayURL,
	)
	return err
}

func (r *OrdersRepo) GetByPublicID(ctx context.Context, publicID string) (*models.Order, error) {
	row := r.DB.QueryRowContext(ctx, `
SELECT
  id, public_id, status, created_at, updated_at,
  source_chain, source_asset, amount_in, target_chain, target_asset, target_address,
  estimated_out, min_out, quote_expiry_at,
  pay_window_seconds,
  mixin_opponent_id, mixin_asset_id, mixin_pay_memo, mixin_pay_url,
  deposit_txid, deposit_tx_detected_at, deposit_credited_at, amount_credited, refund_to_address,
  final_out, swap_ref, withdraw_txid, refund_txid
FROM orders
WHERE public_id = ?
LIMIT 1
`, publicID)

	var o models.Order
	var status string
	var createdAt, updatedAt string
	var quoteExpiry sql.NullString
	var depositTxID, depositDetectedAt, depositCreditedAt sql.NullString
	var amountCredited, refundToAddress sql.NullString
	var finalOut, swapRef, withdrawTxID, refundTxID sql.NullString

	if err := row.Scan(
		&o.ID, &o.PublicID, &status, &createdAt, &updatedAt,
		&o.SourceChain, &o.SourceAsset, &o.AmountIn, &o.TargetChain, &o.TargetAsset, &o.TargetAddress,
		&o.EstimatedOut, &o.MinOut, &quoteExpiry,
		&o.PayWindowSeconds,
		&o.MixinOpponentID, &o.MixinAssetID, &o.MixinPayMemo, &o.MixinPayURL,
		&depositTxID, &depositDetectedAt, &depositCreditedAt, &amountCredited, &refundToAddress,
		&finalOut, &swapRef, &withdrawTxID, &refundTxID,
	); err != nil {
		return nil, err
	}

	o.Status = models.OrderStatus(status)
	ct, _ := time.Parse(time.RFC3339Nano, createdAt)
	ut, _ := time.Parse(time.RFC3339Nano, updatedAt)
	o.CreatedAt = ct
	o.UpdatedAt = ut
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
	if withdrawTxID.Valid {
		o.WithdrawTxID = &withdrawTxID.String
	}
	if refundTxID.Valid {
		o.RefundTxID = &refundTxID.String
	}

	return &o, nil
}

func (r *OrdersRepo) SetDepositCreditedByMemo(ctx context.Context, memo string, snapshotID string, creditedAt time.Time, amountCredited string, assetID string, opponentID string) (int64, error) {
	// Mixin-internal transfer snapshots are already credited to the bot.
	// Match by memo + asset_id. Do NOT require opponent_id (sender) since we may not know it in advance.
	res, err := r.DB.ExecContext(ctx, `
UPDATE orders
SET
  status = ?,
  deposit_txid = COALESCE(deposit_txid, ?),
  deposit_tx_detected_at = COALESCE(deposit_tx_detected_at, ?),
  deposit_credited_at = COALESCE(deposit_credited_at, ?),
  amount_credited = COALESCE(amount_credited, ?),
  refund_to_address = COALESCE(refund_to_address, ?),
  updated_at = ?
WHERE
  status IN (?, ?, ?) AND
  mixin_pay_memo = ? AND
  mixin_asset_id = ?
`,
		string(models.StatusDepositCredited),
		snapshotID,
		creditedAt.Format(time.RFC3339Nano),
		creditedAt.Format(time.RFC3339Nano),
		amountCredited,
		opponentID,
		time.Now().UTC().Format(time.RFC3339Nano),
		string(models.StatusAwaitingDeposit),
		string(models.StatusDepositDetected),
		string(models.StatusDepositPendingMixin),
		memo,
		assetID,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func nullableTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format(time.RFC3339Nano)
}
