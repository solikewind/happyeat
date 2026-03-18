#!/usr/bin/env bash

set -Eeuo pipefail

compose_cmd() {
  if command -v docker-compose >/dev/null 2>&1; then
    docker-compose "$@"
  else
    docker compose "$@"
  fi
}

# 基础文件检查
if [ ! -f docker-compose-prod.yml ]; then
  echo "Error: docker-compose-prod.yml not found."
  exit 1
fi

# 检查数据库是否在线
if ! compose_cmd -f docker-compose-prod.yml ps -q postgres | grep -q .; then
  echo "Error: postgres is not running. Deploy first."
  exit 1
fi

echo "Running database migration..."

# --- 关键修正：去掉 --no-deps 确保网络通畅，并统一指向 /run/config ---
compose_cmd -f docker-compose-prod.yml run --rm happyeat-api /bin/sh -c '
  if [ -f /run/config/happyeatservice.remote.yaml ]; then
    echo "Using remote config for migration..."
    /app/migrate -f /run/config/happyeatservice.remote.yaml
  else
    echo "Using default config for migration..."
    /app/migrate -f /app/etc/happyeatservice.yaml
  fi
'

echo "Migration finished."