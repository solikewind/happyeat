.PHONY: api
# 统一入口：一次生成 menu+table 的 types、handler、logic，避免相互覆盖
api:
	goctl api go --api app/api/v1/central.api --dir app

