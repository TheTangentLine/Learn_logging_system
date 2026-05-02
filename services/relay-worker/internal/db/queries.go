package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OutboxRow struct {
	ID      uuid.UUID
	LogID   uuid.UUID
	Payload []byte
}

// FetchPending selects up to limit PENDING rows inside the provided transaction.
func FetchPending(ctx context.Context, tx pgx.Tx, limit int) ([]OutboxRow, error) {
	rows, err := tx.Query(ctx, fetchPendingSQL, limit)
	if err != nil {
		return nil, fmt.Errorf("fetch pending outbox rows: %w", err)
	}
	defer rows.Close()

	var result []OutboxRow
	for rows.Next() {
		var r OutboxRow
		if err := rows.Scan(&r.ID, &r.LogID, &r.Payload); err != nil {
			return nil, fmt.Errorf("scan outbox row: %w", err)
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

// MarkDone marks the outbox row as successfully published.
func MarkDone(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	if _, err := tx.Exec(ctx, markDoneSQL, id, time.Now().UTC()); err != nil {
		return fmt.Errorf("mark outbox row done (id=%s): %w", id, err)
	}
	return nil
}

// MarkFailed marks the outbox row as failed (will be retried on next poll).
func MarkFailed(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	if _, err := tx.Exec(ctx, markFailedSQL, id, time.Now().UTC()); err != nil {
		return fmt.Errorf("mark outbox row failed (id=%s): %w", id, err)
	}
	return nil
}
