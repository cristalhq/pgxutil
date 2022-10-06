package pgxutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
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

// InWriteTx runs the given function within a transaction with a given isolation level.
func (db *DB) InWriteTx(ctx context.Context, level pgx.TxIsoLevel, fn func(tx pgx.Tx) error) error {
	return db.inTx(ctx, level, "", fn)
}

// InReadTx runs the given function within a read-only transaction with read commited isolation level.
func (db *DB) InReadTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return db.inTx(ctx, pgx.ReadCommitted, pgx.ReadOnly, fn)
}

func (db *DB) inTx(ctx context.Context, level pgx.TxIsoLevel, access pgx.TxAccessMode,
	fn func(tx pgx.Tx) error) (err error) {
	conn, errAcq := db.pool.Acquire(ctx)
	if errAcq != nil {
		return fmt.Errorf("acquiring connection: %w", errAcq)
	}
	defer conn.Release()

	opts := pgx.TxOptions{
		IsoLevel:   level,
		AccessMode: access,
	}

	tx, errBegin := conn.BeginTx(ctx, opts)
	if errBegin != nil {
		return fmt.Errorf("begin tx: %w", errBegin)
	}

	defer func() {
		errRollback := tx.Rollback(ctx)
		if !(errRollback == nil || errors.Is(errRollback, pgx.ErrTxClosed)) {
			err = errRollback
		}
	}()

	if err := fn(tx); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			return fmt.Errorf("rollback tx: %v (original: %w)", errRollback, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
