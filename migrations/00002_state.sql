-- +goose Up

CREATE TABLE IF NOT EXISTS kv (
  k TEXT PRIMARY KEY,
  v TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS kv;
