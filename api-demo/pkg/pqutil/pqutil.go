package pqutil

import (
	"context"
	"database/sql"
)

// Queryer represents the ability to query SQL databases
type Queryer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Transactioner is an interface that allows beginning transactions for a database connection
type Transactioner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Scanner interface {
	Scan(...interface{}) error
}

type Iter interface {
	Next() bool
	Err() error
}

type ScannerIter interface {
	Scanner
	Iter
}

// HasRowsAffected checks to see if the expected number of rows were affected by the call to db.Exec
func HasRowsAffected(res sql.Result, n int64) (bool, error) {
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == n, nil
}
