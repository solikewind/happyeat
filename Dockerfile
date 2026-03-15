# Build stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /project/happyeat

RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build API and migration binaries
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/main ./app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/migrate ./app/cmd/migrate

# Runtime stage
FROM alpine:3.20

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache tzdata ca-certificates wget \
    && addgroup -S app \
    && adduser -S -G app app

ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /out/main /app/main
COPY --from=builder /out/migrate /app/migrate
COPY --from=builder /project/happyeat/app/etc /app/etc

RUN chmod +x /app/main /app/migrate && chown -R app:app /app

USER app

EXPOSE 8888
CMD ["./main"]
