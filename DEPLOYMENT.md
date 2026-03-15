# HappyEat 后端部署指南（Ubuntu + Docker）

本指南用于 Ubuntu 服务器生产部署，默认你已经安装了 Docker。

## 1. 服务器前置条件

- Ubuntu 20.04+
- Docker Engine 20.10+
- Docker Compose 插件（`docker compose version`）
- 已放行端口：`22`（SSH），`80/443`（若走 Nginx），`${API_PORT}`（若直接暴露 API）

如果没有 compose 插件：

```bash
sudo apt update
sudo apt install -y docker-compose-plugin
```

## 2. 拉取项目

```bash
git clone <你的仓库地址> happyeat
cd happyeat
```

## 3. 配置文件准备

### 3.1 创建 `.env`

```bash
cp .env.example .env
nano .env
```

至少修改以下字段：

```env
TZ=Asia/Shanghai
API_PORT=8888
DB_USER=postgres
DB_PASSWORD=<强密码>
DB_NAME=happyeat
JWT_SECRET=<强随机密钥>
```

生成 JWT 随机密钥：

```bash
openssl rand -base64 32
```

### 3.2 创建生产覆盖配置

```bash
cp app/etc/happyeatservice.remote.yaml.example app/etc/happyeatservice.remote.yaml
nano app/etc/happyeatservice.remote.yaml
```

至少修改：

- `Auth.AccessSecret`
- `SqlConfig.DataSource`（密码必须与 `.env` 的 `DB_PASSWORD` 一致）

## 4. 一键部署

```bash
chmod +x deploy.sh migrate.sh
./deploy.sh
```

`deploy.sh` 会自动完成：

- 校验 compose 和配置文件
- 停止旧容器
- 重建镜像
- 启动服务
- 等待 PostgreSQL 健康
- 自动执行数据库迁移

## 5. 部署验证

```bash
docker compose -f docker-compose-prod.yml ps
curl http://127.0.0.1:8888/health
```

预期返回：

```json
{"status":"ok","service":"happyeat-api"}
```

## 6. 日常运维命令

查看日志：

```bash
docker compose -f docker-compose-prod.yml logs -f happyeat-api
docker compose -f docker-compose-prod.yml logs -f postgres
```

重启：

```bash
docker compose -f docker-compose-prod.yml restart
```

仅执行迁移：

```bash
./migrate.sh
```

停止服务：

```bash
docker compose -f docker-compose-prod.yml down
```

## 7. 安全建议

- 首次部署前务必修改所有默认密码与密钥。
- 生产编排已将 PostgreSQL 设为仅容器内网络访问（不映射宿主机端口）。
- 建议接入 Nginx + HTTPS（Let's Encrypt）。
- 建议每天定时备份数据库。

备份示例：

```bash
mkdir -p backups
docker exec happyeat-postgres pg_dump -U postgres happyeat > backups/happyeat_$(date +%Y%m%d_%H%M%S).sql
```

## 8. 故障排查

1. 容器起不来：

```bash
docker compose -f docker-compose-prod.yml logs --tail=200
```

2. 数据库连接失败：

- 检查 `app/etc/happyeatservice.remote.yaml` 中 `SqlConfig.DataSource`
- 检查 `.env` 中 `DB_PASSWORD`

3. 健康检查失败：

```bash
docker exec happyeat-api wget -qO- http://localhost:8888/health
```
