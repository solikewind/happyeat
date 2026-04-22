#!/usr/bin/env bash

# 开启严格错误模式：一旦出错立即停止，未定义变量报错
set -Eeuo pipefail

# 自动适配 docker-compose 或 docker compose 命令
compose_cmd() {
  if command -v docker-compose >/dev/null 2>&1; then
    docker-compose "$@"
  else
    docker compose "$@"
  fi
}

# 1. 基础环境检查
if [ ! -f docker-compose-prod.yml ]; then
  echo "❌ Error: docker-compose-prod.yml not found in current directory."
  exit 1
fi

# 2. 检查数据库是否在线 (Migration 必须依赖运行中的 DB)
# 优先看当前 compose 项目里的 postgres；若为空，再认固定容器名（避免用别的 yml/目录启动后 ps 对不上）
postgres_container_id() {
  compose_cmd -f docker-compose-prod.yml ps -q postgres 2>/dev/null || true
}

if ! postgres_container_id | grep -q .; then
  if ! docker ps -q -f name=happyeat-postgres -f status=running | grep -q .; then
    echo "❌ Error: postgres is not running (compose service postgres or container happyeat-postgres)."
    echo "   Start DB first, e.g.: docker compose -f docker-compose-prod.yml up -d postgres"
    echo "   Wait until health is healthy, then run: ./migrate.sh"
    exit 1
  fi
  echo "ℹ️ Postgres container happyeat-postgres is running (not listed under this compose ps; continuing)."
fi

# 刚启动时可能仍在 health: starting，避免 migrate 立刻连库失败
if docker ps -q -f name=happyeat-postgres -f status=running | grep -q .; then
  echo "⏳ Waiting for Postgres to accept connections..."
  for i in $(seq 1 30); do
    if docker exec happyeat-postgres pg_isready -U "${DB_USER:-postgres}" -d "${DB_NAME:-happyeat}" >/dev/null 2>&1; then
      echo "✅ Postgres is ready."
      break
    fi
    if [ "$i" -eq 30 ]; then
      echo "❌ Postgres did not become ready in time (pg_isready failed)."
      exit 1
    fi
    sleep 2
  done
fi

echo "🚀 Starting database migration..."

# 3. 执行迁移
# 修正点：去掉 --no-deps (确保能连上 postgres 网络)
# 修正点：显式指向挂载后的配置路径 /run/config/
compose_cmd -f docker-compose-prod.yml run --rm happyeat-api /bin/sh -c '
  TARGET_CONF="/run/config/happyeatservice.remote.yaml"
  DEFAULT_CONF="/app/etc/happyeatservice.yaml"

  if [ -f "$TARGET_CONF" ]; then
    echo "📂 Found remote config at $TARGET_CONF, executing..."
    /app/migrate -f "$TARGET_CONF"
  else
    echo "⚠️ Remote config not found, falling back to $DEFAULT_CONF"
    /app/migrate -f "$DEFAULT_CONF"
  fi
'

echo "✅ Migration finished successfully."