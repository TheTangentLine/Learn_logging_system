#!/usr/bin/env bash
# Usage: bash reset.sh
# Wipes the Postgres data volume and re-runs all migrations from scratch.
# Run from the repository root or the scripts/ directory.
set -euo pipefail

# Resolve the repo root (one level up from scripts/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$ROOT_DIR"

echo "Stopping api-server and migrate services..."
docker compose stop api-server migrate 2>/dev/null || true

echo "Tearing down postgres and removing its data volume..."
docker compose down postgres -v

echo "Starting postgres (waiting for healthy)..."
docker compose up postgres --wait

echo "Running migrations..."
docker compose up migrate --wait

echo ""
echo "Database reset complete."
echo "Run 'docker compose up api-server' to bring the API back up."
