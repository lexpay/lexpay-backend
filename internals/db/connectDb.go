package db

import (
	"log/slog"
	"time"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(DatabaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(DatabaseURL)
    if err != nil {
		slog.Error("Failed to parse database URL", slog.String("error", err.Error()))
		return nil, err
	}

	config.MaxConnIdleTime = 5 * time.Minute
	config.MaxConnLifetime = 10 * time.Minute
	config.MaxConns = 100
	config.MinConns = 10

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("Failed to create database connection pool", slog.String("error", err.Error()))
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("Failed to ping database", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("Database connection pool created successfully")

	return pool, nil
}