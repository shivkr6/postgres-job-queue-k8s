package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const insertJobSQL = `INSERT INTO jobs (payload) VALUES ($1::jsonb)`

type DB struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, databaseURL string) (*DB, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *DB) Migrate(ctx context.Context, migrationSQL string) error {
	_, err := db.pool.Exec(ctx, migrationSQL)
	return err
}

func (db *DB) Enqueue(ctx context.Context, payload json.RawMessage) error {
	if !json.Valid(payload) {
		return fmt.Errorf("payload must be valid JSON")
	}

	_, err := db.pool.Exec(ctx, insertJobSQL, string(payload))
	return err
}

func (db *DB) Seed(ctx context.Context, count int) error {
	if count <= 0 {
		return fmt.Errorf("count must be greater than 0")
	}

	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for i := 1; i <= count; i++ {
		payload, err := json.Marshal(struct {
			Type  string `json:"type"`
			Index int    `json:"index"`
		}{
			Type:  "seed",
			Index: i,
		})
		if err != nil {
			return fmt.Errorf("build seed payload %d: %w", i, err)
		}

		_, err = tx.Exec(ctx, insertJobSQL, string(payload))
		if err != nil {
			return fmt.Errorf("enqueue seed job %d: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}

func (db *DB) Close() {
	db.pool.Close()
}
