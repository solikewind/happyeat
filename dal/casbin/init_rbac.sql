-- IAM RBAC 主存表（单租户版本）
-- 说明：权限关系以业务表为主存，casbin_rule 仅作为投影策略表。

CREATE TABLE IF NOT EXISTS iam_users (
  id bigint PRIMARY KEY,
  user_code text NOT NULL UNIQUE,
  display_name text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  delete_ts bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS iam_roles (
  id bigint PRIMARY KEY,
  role_code text NOT NULL UNIQUE,
  role_name text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  delete_ts bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS iam_permissions (
  id bigint PRIMARY KEY,
  permission_code text NOT NULL UNIQUE,
  description text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  delete_ts bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS iam_user_roles (
  iam_user_id bigint NOT NULL REFERENCES iam_users(id) ON DELETE CASCADE,
  iam_role_id bigint NOT NULL REFERENCES iam_roles(id) ON DELETE CASCADE,
  PRIMARY KEY (iam_user_id, iam_role_id)
);

CREATE TABLE IF NOT EXISTS iam_role_permissions (
  iam_role_id bigint NOT NULL REFERENCES iam_roles(id) ON DELETE CASCADE,
  iam_permission_id bigint NOT NULL REFERENCES iam_permissions(id) ON DELETE CASCADE,
  PRIMARY KEY (iam_role_id, iam_permission_id)
);
