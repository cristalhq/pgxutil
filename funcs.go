package pgxutil

import (
	"context"
	"reflect"

	"github.com/cristalhq/builq"
	"github.com/jackc/pgx/v5"
)

type Scannable interface {
	Scan(pgx.Row) error
}

func Exec(ctx context.Context, db *DB, b builq.Builder) error {
	query, args, err := b.Build()
	if err != nil {
		return err
	}

	return db.InWriteTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, query, args...)
		return err
	})
}

func ExecRead[T Scannable](ctx context.Context, db *DB, b builq.Builder, dst T) error {
	query, args, err := b.Build()
	if err != nil {
		return err
	}

	return db.InWriteTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, query, args...)
		return dst.Scan(row)
	})
}

func Read[T Scannable](ctx context.Context, db *DB, b builq.Builder, dst T) error {
	query, args, err := b.Build()
	if err != nil {
		return err
	}

	return db.InReadTx(ctx, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, query, args...)
		return dst.Scan(row)
	})
}

func ReadMany[T Scannable](ctx context.Context, db *DB, b builq.Builder, dst *[]T) error {
	query, args, err := b.Build()
	if err != nil {
		return err
	}

	var v T
	new := makeNew(v)

	return db.InReadTx(ctx, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			row := new().(T)
			if err := row.Scan(rows); err != nil {
				return err
			}
			*dst = append(*dst, row)
		}
		return rows.Err()
	})
}

// Hack to new(...) generic type. There might be better solution.
func makeNew[T any](v T) func() any {
	if typ := reflect.TypeOf(v); typ.Kind() == reflect.Ptr {
		elem := typ.Elem()
		return func() any {
			return reflect.New(elem).Interface()
		}
	}
	return func() any { return new(T) }
}
