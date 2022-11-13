package pgxutil

import "github.com/jackc/pgx/v5"

type Int64 int64

func (s *Int64) Scan(row pgx.Row) error {
	return row.Scan(&s)
}

type String string

func (s *String) Scan(row pgx.Row) error {
	return row.Scan(&s)
}
