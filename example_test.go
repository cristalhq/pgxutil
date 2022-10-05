package pgxutil_test

import (
	"context"

	"github.com/cristalhq/pgxutil"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Example() {
	var pool *pgxpool.Pool

	db, err := pgxutil.New(pool)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// to make transaction with a given isolation level
	level := pgx.Serializable
	errTx := db.InWriteTx(ctx, level, func(tx pgx.Tx) error {
		// TODO: good query with tx
		return nil
	})
	if errTx != nil {
		panic(errTx)
	}

	// to make read-only transaction with a read committed isolation level
	errRead := db.InReadTx(ctx, func(tx pgx.Tx) error {
		// TODO: good read-only query with tx
		return nil
	})
	if errRead != nil {
		panic(errRead)
	}
}
