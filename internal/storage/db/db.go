package db

import (
	"context"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDBPool(ctx context.Context, connString string) (*DB, error) {
	connPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return &DB{pool: connPool}, nil
}

func (db *DB) Ping() error {
	err := db.pool.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}
