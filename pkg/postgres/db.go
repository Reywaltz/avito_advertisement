package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
)

func NewDB(connStr string) (*pgx.Conn, error) {
	if connStr == "" {
		return nil, errors.New("Conn string is not set")
	}

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
