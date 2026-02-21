-- +goose Up

ALTER TABLE orders ADD COLUMN refund_asset_id TEXT;
ALTER TABLE orders ADD COLUMN refund_amount TEXT;
ALTER TABLE orders ADD COLUMN refund_received_snapshot_id TEXT;

CREATE INDEX IF NOT EXISTS idx_orders_refund_received_snapshot_id ON orders(refund_received_snapshot_id);

-- +goose Down

DROP INDEX IF EXISTS idx_orders_refund_received_snapshot_id;
-- SQLite doesn't support DROP COLUMN reliably; leave columns in place.
