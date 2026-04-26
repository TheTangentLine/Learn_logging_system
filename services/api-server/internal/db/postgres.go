package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewWritePool connects to the primary Postgres instance for transactional writes.
func NewWritePool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return newPool(ctx, dsn, "write")
}

// NewReadPool connects to a Postgres read replica for outbox polling and read queries.
func NewReadPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return newPool(ctx, dsn, "read")
}

func newPool(ctx context.Context, dsn, label string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create postgres %s pool: %w", label, err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres %s pool: %w", label, err)
	}
	return pool, nil
}
