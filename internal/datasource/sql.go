package datasource

import (
	"database/sql"
	"errors"
)

type SQL interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

var (
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrNoRows         = errors.New("no rows found")
)
