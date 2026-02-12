.PHONY: api migrate run
# 统一入口：一次生成 menu+table 的 types、handler、logic，避免相互覆盖
api:
	goctl api go --api app/api/v1/central.api --dir app

# 数据库迁移（在 app 目录执行，需先配置 app/etc/happyeatservice.yaml 的 DataSource）
migrate:
	cd app && go run ./cmd/migrate -f etc/happyeatservice.yaml

# 启动 HTTP 服务（需先执行 make migrate，可选执行 dal/casbin/init_policy.sql）
run:
	cd app && go run . -f etc/happyeatservice.yaml

