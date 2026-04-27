CREATE TABLE IF NOT EXISTS outbox (
    id           UUID        PRIMARY KEY,
    log_id       UUID        NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    status       VARCHAR(10) NOT NULL DEFAULT 'PENDING'
                             CHECK (status IN ('PENDING','DONE','FAILED')),
    payload      JSONB       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

-- Partial index: only indexes PENDING rows, stays small as rows transition to DONE/FAILED.
-- Covers the relay worker's: SELECT ... WHERE status = 'PENDING' ORDER BY created_at FOR UPDATE SKIP LOCKED
CREATE INDEX IF NOT EXISTS idx_outbox_pending ON outbox (created_at ASC) WHERE status = 'PENDING';

CREATE INDEX IF NOT EXISTS idx_outbox_log_id  ON outbox (log_id);
