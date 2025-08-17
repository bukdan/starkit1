package db

import (
	"context"
	"time"

	"user-service/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Ping
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	// optional: set pool config options
	pool.Config().MaxConns = 10

	return pool, nil
}
