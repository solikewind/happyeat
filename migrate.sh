#!/usr/bin/env bash

set -Eeuo pipefail

compose_cmd() {
  if command -v docker-compose >/dev/null 2>&1; then
    docker-compose "$@"
  else
    docker compose "$@"
  fi
}

if [ ! -f docker-compose-prod.yml ]; then
  echo "Error: docker-compose-prod.yml not found."
  exit 1
fi

if [ ! -f app/etc/happyeatservice.remote.yaml ]; then
  echo "Error: app/etc/happyeatservice.remote.yaml not found."
  exit 1
fi

if ! compose_cmd -f docker-compose-prod.yml ps -q happyeat-api | grep -q .; then
  echo "Error: happyeat-api is not running. Deploy first."
  exit 1
fi

echo "Running database migration..."
docker exec happyeat-api /bin/sh -c '/app/migrate -f /app/etc/happyeatservice.yaml -rf /run/config/happyeatservice.remote.yaml'
echo "Migration finished."
