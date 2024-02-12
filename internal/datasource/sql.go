package datasource

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type SQL interface {
	Queryer
	Execer
	Transaction
}

type Queryer interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type TxExecer interface {
	Execer
	Queryer
	Rollback() error
	Commit() error
}

type Transaction interface {
	MustBegin() *sqlx.Tx
}

var (
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrNoRows         = errors.New("no rows found")
)
