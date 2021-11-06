package pgxutil

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

var readOnlyOpts = pgx.TxOptions{
	IsoLevel:       pgx.ReadCommitted,
	AccessMode:     pgx.ReadOnly,
	DeferrableMode: "",
}

// New creates a new wrapper for pgx.
func New(pool *pgxpool.Pool) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	db := &DB{
		pool: pool,
	}
	return db, nil
}

// InTx runs the given function within a transaction with a given isolation level.
func (db *DB) InTx(ctx context.Context, level pgx.TxIsoLevel, fn func(tx pgx.Tx) error) error {
	return db.inTx(ctx, level, "", fn)
}

// InReadOnlyTx runs the given function within a read-only transaction with read commited isolation level.
func (db *DB) InReadOnlyTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return db.inTx(ctx, pgx.ReadCommitted, pgx.ReadOnly, fn)
}

func (db *DB) inTx(ctx context.Context, level pgx.TxIsoLevel, access pgx.TxAccessMode, fn func(tx pgx.Tx) error) error {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	opts := pgx.TxOptions{
		IsoLevel:   level,
		AccessMode: access,
	}
	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if err1 := tx.Rollback(ctx); err1 != nil {
			return fmt.Errorf("rolling back transaction: %v (original error: %w)", err1, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
