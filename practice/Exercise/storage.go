package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Connection pool config
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresStore{db: db}, nil
}

func (store *PostgresStore) SaveResults(ctx context.Context, results []*URLCheckResult) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO url_checks
			(url, status, status_code, response_time_ms, error)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, result := range results {
		_, err := stmt.Exec(
			result.URL,
			result.Status,
			result.StatusCode,
			result.ResponseTimeMs,
			result.Error,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
