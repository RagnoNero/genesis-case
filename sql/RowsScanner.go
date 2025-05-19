package sql

type RowsScanner interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}
