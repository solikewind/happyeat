# HappyEat API 部署文档

## 服务器要求

- Ubuntu 20.04+ 或其他 Linux 发行版
- Docker 20.10+
- Docker Compose 2.0+
- 至少 2GB RAM
- 至少 10GB 可用磁盘空间

## 部署步骤

### 1. 在服务器上安装 Docker

```bash
# 更新包索引
sudo apt update

# 安装必要的依赖
sudo apt install -y curl git

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 将当前用户添加到 docker 组（避免每次使用 sudo）
sudo usermod -aG docker $USER
# 重新登录或运行以下命令使组权限生效
newgrp docker

# 安装 Docker Compose
sudo apt install -y docker-compose-plugin
# 或者使用以下命令安装独立版本
# sudo curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-linux-x86_64" -o /usr/local/bin/docker-compose
# sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker compose version
```

### 2. 克隆代码仓库

```bash
# 替换为你的 Git 仓库地址
git clone https://github.com/yourusername/happyeat.git
cd happyeat
```

### 3. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填入正确的配置
nano .env
```

**必须配置的变量：**

```env
# 数据库配置
DB_USER=postgres              # 数据库用户名
DB_PASSWORD=your_strong_password  # 强密码（生产环境必须修改）
DB_NAME=happyeat              # 数据库名称
DB_PORT=5432                 # 数据库端口

# API 服务端口
API_PORT=8888                # 对外暴露的端口

# JWT 配置
JWT_SECRET=your_jwt_secret_change_this  # JWT 密钥（生产环境必须修改为强随机字符串）
JWT_EXPIRE=86400             # Token 过期时间（秒）

# LLM 配置（可选）
LLM_API_KEY=your_api_key     # 如果需要使用 LLM 功能
LLM_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1
```

**生成安全的 JWT Secret：**

```bash
# 生成一个安全的随机字符串
openssl rand -base64 32
```

### 4. 使用部署脚本部署

```bash
# 给部署脚本添加执行权限
chmod +x deploy.sh

# 运行部署脚本
.bash deploy.sh
```

部署脚本会自动完成以下操作：
- 检查环境
- 停止旧服务（如果存在）
- 构建 Docker 镜像
- 启动服务
- 显示服务状态

### 5. 验证部署

```bash
# 检查服务状态
docker compose -f docker-compose-prod.yml ps

# 查看日志
docker compose -f docker-compose-prod.yml logs -f

# 健康检查
curl http://localhost:8888/health
```

期望的响应：

```json
{
  "status": "ok",
  "service": "happyeat-api"
}
```

### 6. 数据库迁移（如果需要）

如果你的项目需要运行数据库迁移，可以在容器启动后运行：

```bash
# 进入 API 容器
docker exec -it happyeat-api /bin/sh

# 运行迁移命令（根据你的项目实际情况调整）
# ./migrate

# 退出容器
exit
```

## 常用命令

### 查看服务状态

```bash
docker compose -f docker-compose-prod.yml ps
```

### 查看日志

```bash
# 查看所有服务日志
docker compose -f docker-compose-prod.yml logs -f

# 只查看 API 服务日志
docker compose -f docker-compose-prod.yml logs -f happyeat-api

# 只查看数据库日志
docker compose -f docker-compose-prod.yml logs -f postgres
```

### 重启服务

```bash
# 重启所有服务
docker compose -f docker-compose-prod.yml restart

# 重启 API 服务
docker compose -f docker-compose-prod.yml restart happyeat-api
```

### 停止服务

```bash
# 停止服务但保留数据
docker compose -f docker-compose-prod.yml down

# 停止服务并删除数据卷（慎用！）
docker compose -f docker-compose-prod.yml down -v
```

### 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建并启动
.bash deploy.sh
```

### 备份数据库

```bash
# 创建备份目录
mkdir -p backups

# 备份数据库
docker exec happyeat-postgres pg_dump -U postgres happyeat > backups/happyeat_$(date +%Y%m%d_%H%M%S).sql
```

### 恢复数据库

```bash
# 从备份文件恢复
docker exec -i happyeat-postgres psql -U postgres happyeat < backups/backup_file.sql
```

## 防火墙配置

如果你的服务器启用了防火墙，需要开放 API 端口：

```bash
# 使用 UFW（Ubuntu 防火墙）
sudo ufw allow 8888/tcp

# 或者使用 iptables
sudo iptables -A INPUT -p tcp --dport 8888 -j ACCEPT
```

## Nginx 反向代理配置（可选）

如果你想使用 Nginx 作为反向代理，可以创建以下配置：

```nginx
# /etc/nginx/sites-available/happyeat
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8888;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/happyeat /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 监控和维护

### 设置日志轮转

```bash
# 创建日志轮转配置
sudo nano /etc/logrotate.d/happyeat
```

内容：

```
/home/yourusername/happyeat/docker-compose-prod.yml {
    daily
    rotate 7
    compress
    missingok
    notifempty
    create 644 yourusername yourusername
    sharedscripts
    postrotate
        docker compose -f /home/yourusername/happyeat/docker-compose-prod.yml restart
    endscript
}
```

### 设置自动备份

创建 cron 任务：

```bash
crontab -e
```

添加以下内容（每天凌晨 2 点备份）：

```
0 2 * * * /home/yourusername/backup-happyeat.sh
```

创建备份脚本：

```bash
#!/bin/bash
# /home/yourusername/backup-happyeat.sh

BACKUP_DIR="/home/yourusername/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# 备份数据库
docker exec happyeat-postgres pg_dump -U postgres happyeat > $BACKUP_DIR/happyeat_$DATE.sql

# 删除 7 天前的备份
find $BACKUP_DIR -name "happyeat_*.sql" -mtime +7 -delete
```

## 故障排查

### 服务无法启动

1. 检查日志：
```bash
docker compose -f docker-compose-prod.yml logs
```

2. 检查环境变量：
```bash
cat .env
```

3. 检查端口占用：
```bash
sudo lsof -i :8888
```

### 数据库连接失败

1. 检查数据库容器状态：
```bash
docker ps | grep postgres
```

2. 测试数据库连接：
```bash
docker exec happyeat-postgres psql -U postgres -d happyeat -c "SELECT 1"
```

### 健康检查失败

1. 检查 API 容器内服务是否正常：
```bash
docker exec happyeat-api wget -O- http://localhost:8888/health
```

2. 检查容器资源使用情况：
```bash
docker stats
```

## 安全建议

1. **修改默认密码**：在生产环境中必须修改所有默认密码
2. **使用强 JWT Secret**：生成足够长的随机字符串
3. **限制数据库端口暴露**：生产环境中不要将数据库端口暴露到主机网络
4. **定期更新**：定期更新 Docker 镜像和依赖
5. **启用 HTTPS**：使用 Nginx + Certbot 启用 HTTPS
6. **备份策略**：定期备份数据库和配置文件
7. **监控日志**：设置日志监控，及时发现异常

## 联系支持

如有问题，请查看：
- 项目文档：[README.md](README.md)
- Issue 跟踪：https://github.com/yourusername/happyeat/issues
