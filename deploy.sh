#!/usr/bin/env bash

set -Eeuo pipefail

compose_cmd() {
  if command -v docker-compose >/dev/null 2>&1; then
    docker-compose "$@"
  else
    docker compose "$@"
  fi
}

run_migration() {
  compose_cmd -f docker-compose-prod.yml run --rm --no-deps happyeat-api /bin/sh -c '
    if /app/migrate -h 2>&1 | grep -q -- " -rf "; then
      if [ -f /run/config/happyeatservice.remote.yaml ]; then
        /app/migrate -f /app/etc/happyeatservice.yaml -rf /run/config/happyeatservice.remote.yaml
      else
        /app/migrate -f /app/etc/happyeatservice.yaml
      fi
    else
      /app/migrate -f /app/etc/happyeatservice.yaml
    fi
  '
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

echo "[1/7] Checking compose file..."
compose_cmd -f docker-compose-prod.yml config >/dev/null

echo "[2/7] Stopping old services..."
compose_cmd -f docker-compose-prod.yml down || true

echo "[3/7] Building images..."
compose_cmd -f docker-compose-prod.yml build --pull

echo "[4/7] Starting database..."
compose_cmd -f docker-compose-prod.yml up -d postgres

echo "[5/7] Waiting for database to be ready..."
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

if [ "${RUN_MIGRATION:-0}" = "1" ]; then
  echo "[6/7] Running migration..."
  run_migration
else
  echo "[6/7] Skip migration (set RUN_MIGRATION=1 to enable)"
fi

echo "[7/7] Starting API service..."
compose_cmd -f docker-compose-prod.yml up -d happyeat-api

echo "Waiting for API health..."
for i in {1..60}; do
  state="$(docker inspect -f '{{.State.Status}}' happyeat-api 2>/dev/null || true)"
  health="$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' happyeat-api 2>/dev/null || true)"
  if [ "$state" != "running" ]; then
    echo "API container is not running (state: ${state:-unknown}). Recent logs:"
    compose_cmd -f docker-compose-prod.yml logs --tail=120 happyeat-api || true
    exit 1
  fi
  if [ "$health" = "healthy" ]; then
    break
  fi
  if [ "$health" = "unhealthy" ]; then
    echo "API failed to become healthy. Recent logs:"
    compose_cmd -f docker-compose-prod.yml logs --tail=120 happyeat-api || true
    exit 1
  fi
  if [ "$i" -eq 60 ]; then
    echo "Error: API did not become healthy in time."
    compose_cmd -f docker-compose-prod.yml logs --tail=120 happyeat-api || true
    exit 1
  fi
  sleep 1
done

echo "Service status:"
compose_cmd -f docker-compose-prod.yml ps

echo "Done. Health check: http://localhost:${API_PORT:-8888}/health"
