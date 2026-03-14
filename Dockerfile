# 第一阶段：编译 (使用 2026 年最新的 1.25.6 版本)
FROM golang:1.25.6-alpine AS builder

WORKDIR /project/happyeat

# 先拷贝依赖文件（利用 Docker 缓存层，加速后续构建）
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 拷贝源码
COPY . .

# 编译优化：
# CGO_ENABLED=0 确保静态链接，不需要依赖宿主机的 C 库
# -ldflags="-s -w" 进一步压缩体积
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# 第二阶段：运行 (使用极简镜像)
FROM alpine:3.20

# 合并换源与安装，减少镜像层数
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \    
    apk add --no-cache tzdata ca-certificates

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 拷贝二进制文件
COPY --from=builder /project/happyeat/main .
# 如果有配置目录，再单独考进来
# COPY --from=builder /app/config ./config

# 如果你的程序需要读取配置文件（如 config.yaml）
# COPY config.yaml .

EXPOSE 8080

CMD ["./main"]