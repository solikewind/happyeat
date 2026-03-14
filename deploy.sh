#!/bin/bash

# HappyEat API 部署脚本
# 使用方法: bash deploy.sh

set -e  # 遇到错误立即退出

echo "=========================================="
echo "  HappyEat API 部署脚本"
echo "=========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    echo "请先安装 Docker: https://docs.docker.com/engine/install/"
    exit 1
fi

# 检查 Docker Compose 是否安装（支持 docker-compose 或 docker compose）
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
    echo -e "${GREEN}✓ 使用 docker-compose${NC}"
elif docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
    echo -e "${GREEN}✓ 使用 docker compose（插件版本）${NC}"
else
    echo -e "${RED}错误: Docker Compose 未安装${NC}"
    echo "请先安装 Docker Compose: https://docs.docker.com/compose/install/"
    echo "或者安装 Docker Compose 插件:"
    echo "  sudo apt install -y docker-compose-plugin"
    exit 1
fi

# 检查 .env 文件是否存在
if [ ! -f .env ]; then
    echo -e "${YELLOW}警告: .env 文件不存在${NC}"
    if [ -f .env.example ]; then
        echo "从 .env.example 创建 .env 文件..."
        cp .env.example .env
        echo -e "${YELLOW}请编辑 .env 文件，填入正确的配置后重新运行部署脚本${NC}"
        exit 1
    else
        echo -e "${RED}错误: .env.example 文件也不存在${NC}"
        exit 1
    fi
fi

echo -e "${GREEN}✓ 环境检查通过${NC}"

# 停止旧服务（如果存在）
echo ""
echo "停止旧服务..."
if $DOCKER_COMPOSE -f docker-compose-prod.yml ps -q 2>/dev/null | grep -q .; then
    $DOCKER_COMPOSE -f docker-compose-prod.yml down
    echo -e "${GREEN}✓ 旧服务已停止${NC}"
else
    echo "没有运行中的服务"
fi

# 构建新镜像
echo ""
echo "构建 Docker 镜像..."
$DOCKER_COMPOSE -f docker-compose-prod.yml build --no-cache
echo -e "${GREEN}✓ 镜像构建完成${NC}"

# 启动服务
echo ""
echo "启动服务..."
$DOCKER_COMPOSE -f docker-compose-prod.yml up -d
echo -e "${GREEN}✓ 服务已启动${NC}"

# 等待数据库就绪
echo ""
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

# 运行数据库迁移
echo ""
echo "运行数据库迁移..."
if docker exec happyeat-api /app/migrate -f /app/etc/happyeatservice.yaml; then
    echo -e "${GREEN}✓ 数据库迁移完成${NC}"
else
    echo -e "${YELLOW}注意: 数据库迁移失败或已存在${NC}"
    echo -e "${YELLOW}查看日志以获取详细信息${NC}"
fi

# 等待 API 服务就绪
echo ""
echo "等待 API 服务就绪..."
sleep 3
max_wait=30
waited=0
until curl -s http://localhost:8888/health > /dev/null 2>&1; do
    if [ $waited -ge $max_wait ]; then
        echo -e "${YELLOW}警告: API 服务可能未完全就绪${NC}"
        break
    fi
    echo -n "."
    sleep 1
    waited=$((waited + 1))
done
echo ""

# 检查服务状态
echo ""
echo "检查服务状态..."
$DOCKER_COMPOSE -f docker-compose-prod.yml ps

# 显示日志
echo ""
echo "=========================================="
echo "  部署完成！"
echo "=========================================="
echo ""
echo "查看日志: $DOCKER_COMPOSE -f docker-compose-prod.yml logs -f"
echo "停止服务: $DOCKER_COMPOSE -f docker-compose-prod.yml down"
echo "重启服务: $DOCKER_COMPOSE -f docker-compose-prod.yml restart"
echo ""
echo "数据库相关:"
echo "  查看表: docker exec happyeat-postgres psql -U postgres -d happyeat -c '\dt'"
echo "  进入数据库: docker exec -it happyeat-postgres psql -U postgres -d happyeat"
echo "  备份数据: docker exec happyeat-postgres pg_dump -U postgres happyeat > backup.sql"
echo ""
echo -e "${GREEN}服务访问地址: http://localhost:8888${NC}"
echo -e "${GREEN}健康检查: http://localhost:8888/health${NC}"
echo ""
