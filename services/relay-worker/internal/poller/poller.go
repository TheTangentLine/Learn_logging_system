package poller

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"relay-worker/internal/db"
	"relay-worker/internal/producer"
)

// Poller runs the main LISTEN/NOTIFY + fallback-poll loop.
type Poller struct {
	pool         *pgxpool.Pool
	listenConn   *pgx.Conn
	prod         *producer.Producer
	batchSize    int
	pollInterval time.Duration
}

// New creates a Poller. listenConn must already be listening on "outbox_ready".
func New(pool *pgxpool.Pool, listenConn *pgx.Conn, prod *producer.Producer, batchSize int, pollInterval time.Duration) *Poller {
	return &Poller{
		pool:         pool,
		listenConn:   listenConn,
		prod:         prod,
		batchSize:    batchSize,
		pollInterval: pollInterval,
	}
}

// Run blocks until ctx is cancelled, processing outbox batches on each
// notification or fallback-poll timeout.
func (p *Poller) Run(ctx context.Context) error {
	log.Printf("poller: starting (batch=%d, fallback_interval=%s)", p.batchSize, p.pollInterval)

	for {
		// Block until a pg_notify fires OR the fallback timeout elapses.
		waitCtx, cancel := context.WithTimeout(ctx, p.pollInterval)
		_, err := p.listenConn.WaitForNotification(waitCtx)
		cancel()

		if ctx.Err() != nil {
			// Parent context cancelled — clean shutdown.
			log.Println("poller: context cancelled, stopping")
			return nil
		}

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			// Unexpected error on the listen connection.
			return err
		}

		if err := p.processBatch(ctx); err != nil {
			log.Printf("poller: batch error: %v", err)
			// Non-fatal — log and keep looping rather than crashing.
		}
	}
}

// processBatch fetches up to batchSize PENDING rows inside a single transaction,
// publishes each to RabbitMQ, and marks them DONE or FAILED.
func (p *Poller) processBatch(ctx context.Context) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := db.FetchPending(ctx, tx, p.batchSize)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	log.Printf("poller: processing %d row(s)", len(rows))

	for _, row := range rows {
		pubErr := p.prod.Publish(ctx, row.Payload)
		if pubErr != nil {
			log.Printf("poller: publish failed for outbox id=%s: %v", row.ID, pubErr)
			if markErr := db.MarkFailed(ctx, tx, row.ID); markErr != nil {
				log.Printf("poller: mark failed error: %v", markErr)
			}
			continue
		}
		if markErr := db.MarkDone(ctx, tx, row.ID); markErr != nil {
			log.Printf("poller: mark done error: %v", markErr)
		}
	}

	return tx.Commit(ctx)
}
