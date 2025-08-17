
package main

import (
    "context"
    "log"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func MustInitDB() *pgxpool.Pool {
    dsn := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/usersvc?sslmode=disable")
    cfg, err := pgxpool.ParseConfig(dsn)
    if err != nil { log.Fatal(err) }
    cfg.MaxConns = 5
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    pool, err := pgxpool.NewWithConfig(ctx, cfg)
    if err != nil { log.Fatal(err) }
    if err := pool.Ping(ctx); err != nil { log.Fatal(err) }
    return pool
}
