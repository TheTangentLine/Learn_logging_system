package db

const insertLog = `
INSERT INTO logs (id, level, message, service_name, timestamp, metadata, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)`

const insertOutbox = `
INSERT INTO outbox (id, log_id, status, payload, created_at)
VALUES ($1, $2, 'PENDING', $3, NOW())`
