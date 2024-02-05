package db

import (
	"context"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	db *pgxpool.Pool
}

func NewDBPool(ctx context.Context, connString string) (*DB, error) {
	connPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return &DB{db: connPool}, nil
}

func (db *DB) Ping() error {
	return db.Ping()
}
