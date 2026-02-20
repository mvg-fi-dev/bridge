-- +goose Up

ALTER TABLE orders ADD COLUMN exinswap_trace_id TEXT;
CREATE INDEX IF NOT EXISTS idx_orders_exinswap_trace_id ON orders(exinswap_trace_id);

-- +goose Down

DROP INDEX IF EXISTS idx_orders_exinswap_trace_id;
-- SQLite doesn't support DROP COLUMN reliably; leave column in place.
