#!/usr/bin/env bash
# Usage: bash seed.sh
#        API_URL=http://localhost:8080 bash seed.sh
set -euo pipefail

BASE_URL="${API_URL:-http://localhost:8080}"

echo "Checking API health at $BASE_URL ..."
curl -sf "$BASE_URL/health" > /dev/null || {
  echo "ERROR: API not reachable at $BASE_URL"
  echo "Make sure the stack is running: docker compose up"
  exit 1
}
echo "API is healthy. Seeding logs..."

post_log() {
  local description="$1"
  local payload="$2"
  local status
  status=$(curl -sf -o /dev/null -w "%{http_code}" \
    -X POST "$BASE_URL/logs" \
    -H "Content-Type: application/json" \
    -d "$payload")
  echo "  [$status] $description"
}

post_log "INFO  — auth-service started" '{
  "level": "INFO",
  "message": "service started successfully",
  "service_name": "auth-service",
  "timestamp": "2026-04-28T00:00:00Z",
  "metadata": {"version": "1.0.0", "port": 8081}
}'

post_log "DEBUG — auth-service token validated" '{
  "level": "DEBUG",
  "message": "JWT token validated",
  "service_name": "auth-service",
  "timestamp": "2026-04-28T00:01:00Z",
  "metadata": {"user_id": "usr_abc123", "ttl_seconds": 3600}
}'

post_log "WARN  — payment-service slow query" '{
  "level": "WARN",
  "message": "database query exceeded 500ms threshold",
  "service_name": "payment-service",
  "timestamp": "2026-04-28T00:02:00Z",
  "metadata": {"query": "SELECT * FROM transactions", "duration_ms": 732}
}'

post_log "ERROR — payment-service charge failed" '{
  "level": "ERROR",
  "message": "payment charge failed: insufficient funds",
  "service_name": "payment-service",
  "timestamp": "2026-04-28T00:03:00Z",
  "metadata": {"order_id": "ord_xyz789", "amount_cents": 4999, "attempt": 1}
}'

post_log "INFO  — notification-service email sent" '{
  "level": "INFO",
  "message": "order confirmation email dispatched",
  "service_name": "notification-service",
  "timestamp": "2026-04-28T00:04:00Z",
  "metadata": {"recipient": "user@example.com", "template": "order_confirm"}
}'

post_log "FATAL — notification-service queue unreachable" '{
  "level": "FATAL",
  "message": "cannot connect to message queue after 5 retries",
  "service_name": "notification-service",
  "timestamp": "2026-04-28T00:05:00Z",
  "metadata": {"host": "rabbitmq", "port": 5672, "retries": 5}
}'

echo ""
echo "Seeding complete. Query logs with:"
echo "  curl 'http://localhost:8080/logs?q=failed'"
echo "  curl 'http://localhost:8080/logs?level=ERROR'"
echo "  curl 'http://localhost:8080/logs?service=auth-service'"
