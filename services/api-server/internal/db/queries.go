package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-server/internal/model"
)

func InsertLogWithOutbox(ctx context.Context, pool *pgxpool.Pool, l *model.Log) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, insertLog,
		l.ID, l.Level, l.Message, l.ServiceName, l.Timestamp, l.Metadata, l.CreatedAt,
	); err != nil {
		return fmt.Errorf("insert log: %w", err)
	}

	payload, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("marshal outbox payload: %w", err)
	}

	if _, err = tx.Exec(ctx, insertOutbox,
		uuid.New(), l.ID, payload,
	); err != nil {
		return fmt.Errorf("insert outbox: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
