package db

const fetchPendingSQL = `
SELECT id, log_id, payload
FROM outbox
WHERE status = 'PENDING'
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED`

const markDoneSQL = `
UPDATE outbox
SET status = 'DONE', processed_at = $2
WHERE id = $1`

const markFailedSQL = `
UPDATE outbox
SET status = 'FAILED', processed_at = $2
WHERE id = $1`