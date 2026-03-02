.PHONY: api migrate run swagger swagger-serve
# 统一入口：一次生成 menu+table 的 types、handler、logic，避免相互覆盖
api:
	goctl api go --api app/api/v1/central.api --dir app

# 数据库迁移（在 app 目录执行，需先配置 app/etc/happyeatservice.yaml 的 DataSource）
migrate:
	cd app && go run ./cmd/migrate -f etc/happyeatservice.yaml

# 数据库生成（ent 会生成 dal/model/ent 下的 client、hook、intercept 等；--feature intercept 启用拦截器）
# 生成后 runtime.go 会被覆盖，需执行 make fix-runtime 打破 import cycle，否则 intercept 会报 missing metadata
generate:
	cd dal/model/ent && go run entgo.io/ent/cmd/ent generate --feature intercept,schema/snapshot ./schema
#cd dal/model/ent && go run -mod=mod entgo.io/ent/cmd/ent generate --feature intercept ./schema

# 生成后执行：用存根覆盖 runtime.go，打破 import cycle（无 cp 时请手动把 runtime_stub.go 复制为 runtime.go）
fix-runtime:
	@cp dal/model/ent/runtime_stub.go dal/model/ent/runtime.go

# 启动 HTTP 服务（需先执行 make migrate，可选执行 dal/casbin/init_policy.sql）
run:
	cd app && go run . -f etc/happyeatservice.yaml

# 根据 central.api 生成 Swagger 接口文档（JSON 默认生成到项目根目录 happyeat.json；需 goctl >= 1.8.2）
swagger:
	goctl api swagger --api app/api/v1/central.api --dir . --filename happyeat

# 启动静态服务供 Swagger 预览：执行后浏览器打开 https://editor.swagger.io 并填入 http://localhost:3780/happyeat.json
swagger-serve:
	npx --yes http-server . -p 3780 -c-1 --cors 