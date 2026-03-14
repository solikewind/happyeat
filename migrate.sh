#!/bin/bash

# HappyEat API 数据库迁移脚本
# 使用方法: bash migrate.sh

set -e  # 遇到错误立即退出

echo "=========================================="
echo "  HappyEat API 数据库迁移"
echo "=========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Docker Compose 命令
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
elif docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
else
    echo -e "${RED}错误: Docker Compose 未安装${NC}"
    exit 1
fi

# 检查 API 容器是否存在
if ! $DOCKER_COMPOSE -f docker-compose-prod.yml ps -q happyeat-api 2>/dev/null | grep -q .; then
    echo -e "${YELLOW}警告: happyeat-api 容器未运行${NC}"
    echo "请先启动服务: bash deploy.sh"
    echo ""
    echo "或者使用以下命令手动运行迁移:"
    echo "  docker exec happyeat-api /app/migrate -f /app/etc/happyeatservice.yaml"
    exit 1
fi

# 等待数据库就绪
echo "等待数据库就绪..."
max_wait=30
waited=0
while ! docker exec happyeat-postgres pg_isready -U postgres > /dev/null 2>&1; do
    if [ $waited -ge $max_wait ]; then
        echo -e "${RED}错误: 数据库未就绪${NC}"
        exit 1
    fi
    echo -n "."
    sleep 1
    waited=$((waited + 1))
done
echo ""
echo -e "${GREEN}✓ 数据库已就绪${NC}"

# 运行迁移
echo ""
echo "开始数据库迁移..."
if docker exec happyeat-api /app/migrate -f /app/etc/happyeatservice.yaml; then
    echo -e "${GREEN}✓ 数据库迁移完成${NC}"
else
    echo -e "${RED}错误: 数据库迁移失败${NC}"
    exit 1
fi

# 查看表
echo ""
echo "查看数据库表..."
docker exec happyeat-postgres psql -U postgres -d happyeat -c "\dt"

echo ""
echo "=========================================="
echo "  迁移完成！"
echo "=========================================="
