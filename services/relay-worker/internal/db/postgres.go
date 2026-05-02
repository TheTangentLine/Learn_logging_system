package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewWritePool creates a connection pool for SELECT FOR UPDATE and UPDATE operations.
func NewWritePool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create postgres write pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres write pool: %w", err)
	}
	return pool, nil
}

// NewListenConn creates a single dedicated connection for LISTEN/NOTIFY.
// This connection must stay alive for the lifetime of the process.
func NewListenConn(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect postgres listen conn: %w", err)
	}
	if _, err := conn.Exec(ctx, "LISTEN outbox_ready"); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("LISTEN outbox_ready: %w", err)
	}
	return conn, nil
}
