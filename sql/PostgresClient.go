package sql

import (
	"fmt"
	"weather-subscription/config"
)

type PostgresClient struct {
	*SqlClient
}

func NewPostgresClient(cfg config.DatabaseConfiguration) (*PostgresClient, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName,
	)

	sqlClient, err := NewSqlClient("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresClient{
		SqlClient: sqlClient,
	}, nil
}
