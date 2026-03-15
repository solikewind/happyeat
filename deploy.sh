#!/usr/bin/env bash

set -Eeuo pipefail

compose_cmd() {
  if command -v docker-compose >/dev/null 2>&1; then
    docker-compose "$@"
  else
    docker compose "$@"
  fi
}

echo "=========================================="
echo "  HappyEat API Production Deploy"
echo "=========================================="

if ! command -v docker >/dev/null 2>&1; then
  echo "Error: Docker is not installed."
  exit 1
fi

if [ ! -f docker-compose-prod.yml ]; then
  echo "Error: docker-compose-prod.yml not found."
  exit 1
fi

if [ ! -f .env ]; then
  echo "Error: .env not found. Please copy from .env.example first."
  exit 1
fi

# Export .env so shell checks use the same values as compose.
set -a
. ./.env
set +a

if [ ! -f app/etc/happyeatservice.remote.yaml ]; then
  echo "Error: app/etc/happyeatservice.remote.yaml not found."
  echo "Hint: cp app/etc/happyeatservice.remote.yaml.example app/etc/happyeatservice.remote.yaml"
  exit 1
fi

if [ "${DB_PASSWORD:-}" = "change-this-password" ]; then
  echo "Error: DB_PASSWORD is still placeholder in .env."
  exit 1
fi

if grep -q "replace-with-strong-jwt-secret\|replace-db-password" app/etc/happyeatservice.remote.yaml; then
  echo "Error: app/etc/happyeatservice.remote.yaml still contains placeholder values."
  exit 1
fi

echo "[1/6] Checking compose file..."
compose_cmd -f docker-compose-prod.yml config >/dev/null

echo "[2/6] Stopping old services..."
compose_cmd -f docker-compose-prod.yml down || true

echo "[3/6] Building images..."
compose_cmd -f docker-compose-prod.yml build --pull

echo "[4/6] Starting services..."
compose_cmd -f docker-compose-prod.yml up -d

echo "[5/6] Waiting for database to be ready..."
for i in {1..60}; do
  if docker exec happyeat-postgres /bin/sh -c 'pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"' >/dev/null 2>&1; then
    break
  fi
  if [ "$i" -eq 60 ]; then
    echo "Error: database is not ready after 60 seconds."
    exit 1
  fi
  sleep 1
done

echo "[6/6] Running migration..."
docker exec happyeat-api /bin/sh -c 'if /app/migrate -h 2>&1 | grep -q -- " -rf "; then /app/migrate -f /app/etc/happyeatservice.yaml -rf /run/config/happyeatservice.remote.yaml; else /app/migrate -f /app/etc/happyeatservice.yaml; fi'

echo "Service status:"
compose_cmd -f docker-compose-prod.yml ps

echo "Done. Health check: http://localhost:${API_PORT:-8888}/health"
