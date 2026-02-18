package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	SQL *sql.DB
}

func Open(sqlitePath string) (*DB, error) {
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_busy_timeout=5000", sqlitePath)
	s, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err := s.Ping(); err != nil {
		return nil, err
	}
	return &DB{SQL: s}, nil
}
