-- Casbin 策略表：与 pckhoi/casbin-pgx-adapter 建表结构一致（id 为 text 主键，列为 p_type 与 v0~v5）。
-- 若使用 adapter 启动，表会自动创建；也可本脚本手动建表并插入示例数据。

CREATE TABLE IF NOT EXISTS casbin_rule (
  id     text PRIMARY KEY,
  p_type text,
  v0     text,
  v1     text,
  v2     text,
  v3     text,
  v4     text,
  v5     text
);

-- 示例（单租户 RBAC）：user_id=dev-admin 继承 super_admin，并允许访问权限矩阵接口。
-- 注意：当前 matcher 为严格 obj/act 匹配，不建议使用 '*' 作为通配符。
INSERT INTO casbin_rule (id, p_type, v0, v1, v2) VALUES ('p-super-admin-role-permissions-get', 'p', 'super_admin', '/central/v1/rbac/role-permissions', 'GET');
INSERT INTO casbin_rule (id, p_type, v0, v1) VALUES ('g-dev-admin-super-admin', 'g', 'dev-admin', 'super_admin');
