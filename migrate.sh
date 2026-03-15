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

if ! compose_cmd -f docker-compose-prod.yml ps -q postgres | grep -q .; then
  echo "Error: postgres is not running. Deploy first."
  exit 1
fi

echo "Running database migration..."
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
echo "Migration finished."
