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
if ! compose_cmd -f docker-compose-prod.yml ps -q postgres | grep -q .; then
  echo "❌ Error: postgres container is not running. Please deploy first."
  exit 1
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