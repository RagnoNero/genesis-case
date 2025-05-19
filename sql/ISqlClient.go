package sql

import "database/sql"

type ISqlClient interface {
	GetDb() sql.DB
	Create(query string, args ...any) (int64, error)
	Read(query string, args ...any) (RowsScanner, error)
	Update(query string, args ...any) (int64, error)
	Delete(query string, args ...any) (int64, error)
	Close() error
}
