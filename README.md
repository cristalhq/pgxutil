# pgxutil

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]

Go [jackc/pgx](https://github.com/jackc/pgx) helper to write proper transactions.


## Features

* Simple API.

## Install

Go version 1.17+

```
go get github.com/cristalhq/pgxutil
```

## Example

```go
// create jackc/pgx pool
var pool *pgxpool.Pool

db, err := pgxutil.New(pool)
if err != nil {
	panic(err)
}

ctx := context.Background()

// to make transaction with a given isolation level
level := pgx.Serializable
errTx := db.InTx(ctx, level, func(tx pgx.Tx) error {
	// TODO: good query with tx
	return nil
})
if errTx != nil {
	panic(errTx)
}

// to make read-only transaction with a read committed isolation level
errRead := db.InReadOnlyTx(ctx, func(tx pgx.Tx) error {
	// TODO: good read-only query with tx
	return nil
})
if errRead != nil {
	panic(errRead)
}	
```

Also see examples: [examples_test.go](https://github.com/cristalhq/pgxutil/blob/main/example_test.go).

## Documentation

See [these docs][pkg-url].

## License

[MIT License](LICENSE).

[build-img]: https://github.com/cristalhq/pgxutil/workflows/build/badge.svg
[build-url]: https://github.com/cristalhq/pgxutil/actions
[pkg-img]: https://pkg.go.dev/badge/cristalhq/pgxutil
[pkg-url]: https://pkg.go.dev/github.com/cristalhq/pgxutil
[reportcard-img]: https://goreportcard.com/badge/cristalhq/pgxutil
[reportcard-url]: https://goreportcard.com/report/cristalhq/pgxutil
[coverage-img]: https://codecov.io/gh/cristalhq/pgxutil/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/cristalhq/pgxutil