package db

import (
	"context"
	"database/sql"
	"time"
)

type StateRepo struct{ DB *sql.DB }

func NewStateRepo(db *sql.DB) *StateRepo { return &StateRepo{DB: db} }

func (r *StateRepo) Get(ctx context.Context, key string) (string, bool, error) {
	row := r.DB.QueryRowContext(ctx, `SELECT v FROM kv WHERE k = ? LIMIT 1`, key)
	var v string
	if err := row.Scan(&v); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return v, true, nil
}

func (r *StateRepo) Set(ctx context.Context, key, value string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := r.DB.ExecContext(ctx, `INSERT INTO kv(k,v,updated_at) VALUES(?,?,?)
ON CONFLICT(k) DO UPDATE SET v=excluded.v, updated_at=excluded.updated_at`, key, value, now)
	return err
}
