package sql

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type SqlClient struct {
	db *sql.DB
}

func NewSqlClient(driver, dsn string) (*SqlClient, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &SqlClient{db: db}, nil
}

func (c *SqlClient) GetDb() *sql.DB {
	return c.db
}

func (c *SqlClient) Create(query string, args ...any) (int64, error) {
	result, err := c.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (c *SqlClient) Read(query string, args ...any) (RowsScanner, error) {
	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *SqlClient) Update(query string, args ...any) (int64, error) {
	result, err := c.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (c *SqlClient) Delete(query string, args ...any) (int64, error) {
	result, err := c.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (c *SqlClient) Close() error {
	if c.db == nil {
		return errors.New("db is not initialized")
	}
	return c.db.Close()
}
