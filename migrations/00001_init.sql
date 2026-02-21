-- +goose Up

CREATE TABLE IF NOT EXISTS orders (
  id TEXT PRIMARY KEY,
  public_id TEXT NOT NULL UNIQUE,
  status TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,

  source_chain TEXT NOT NULL,
  source_asset TEXT NOT NULL,
  amount_in TEXT NOT NULL,
  target_chain TEXT NOT NULL,
  target_asset TEXT NOT NULL,
  target_address TEXT NOT NULL,

  estimated_out TEXT NOT NULL,
  min_out TEXT NOT NULL,
  quote_expiry_at TEXT,

  pay_window_seconds INTEGER NOT NULL,

  mixin_opponent_id TEXT,
  mixin_asset_id TEXT,
  mixin_pay_memo TEXT,
  mixin_pay_url TEXT,

  deposit_txid TEXT,
  deposit_tx_detected_at TEXT,
  deposit_credited_at TEXT,
  amount_credited TEXT,
  refund_to_address TEXT,

  final_out TEXT,
  swap_ref TEXT,
  exinswap_trace_id TEXT,
  withdraw_txid TEXT,
  refund_txid TEXT,
  refund_asset_id TEXT,
  refund_amount TEXT,
  refund_received_snapshot_id TEXT
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);

CREATE TABLE IF NOT EXISTS mixin_snapshots (
  snapshot_id TEXT PRIMARY KEY,
  received_at TEXT NOT NULL,
  raw_json TEXT NOT NULL,

  -- extracted
  created_at TEXT,
  amount TEXT,
  asset_id TEXT,
  opponent_id TEXT,
  memo TEXT
);

CREATE INDEX IF NOT EXISTS idx_mixin_snapshots_created_at ON mixin_snapshots(created_at);

-- +goose Down
DROP TABLE IF EXISTS mixin_snapshots;
DROP TABLE IF EXISTS orders;
