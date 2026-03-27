# HappyEat 实施计划（V1）

> 目标：先把权限与用户体系打稳，再推进业务逻辑优化和前端联调。

## 总体阶段

### 阶段 1：用户与权限模型（当前优先）
- 说明：角色和权限由前端设计，后端负责承接、校验与执行鉴权。
- 产出：
  - 角色定义文档（前端主导）
  - 权限点清单（前端主导，后端确认可映射到 API）
  - 后端 RBAC 数据模型与 Casbin 接入

### 阶段 2：Casbin 鉴权全量接入
- API 默认拒绝，按策略放行。
- 覆盖菜单、餐桌、订单、工作台核心接口。

### 阶段 3：核心业务逻辑优化
- 订单状态机、幂等、并发安全、审计日志。

### 阶段 4：前端联调与权限驱动页面
- 动态路由、按钮级权限、异常场景联调。

---

## 阶段 1 执行细化（从这里开始）

## 1.1 角色与权限定义（前端负责）
- 负责人：前端。
- 后端配合事项：
  - 提供当前 API 列表（path + method + 功能描述）。
  - 对齐权限粒度建议：`module:resource:action`（例如 `order:order:update_status`）。
- 交付物：
  - `角色列表`（如店长、收银、后厨、服务员等）
  - `角色-权限矩阵`
  - `页面/按钮权限点` 到 `API 权限点` 的映射表

### 1.1 完成标准
- 每个页面操作都能映射到一个或多个后端权限点。
- 存在未定义权限的接口为 0。

## 1.2 后端权限模型落地（后端负责）
- 建议最小数据结构：
  - `users`（用户）
  - `roles`（角色）
  - `user_roles`（用户角色关系）
  - `role_policies`（角色权限点，可与 Casbin policy 同步）
- 说明：角色命名和权限点内容以前端设计为准，后端仅做存储与执行。

### 1.2 完成标准
- 能通过用户查询到其角色集合与权限点集合。
- 权限点可稳定映射到 Casbin `sub, obj, act`。

## 1.3 Casbin 接入方案（后端负责）
- 约定：
  - `sub`：用户 ID 或角色 ID（推荐角色）
  - `obj`：资源路径或资源标识（如 `v1/order`）
  - `act`：动作（GET/POST/PUT/DELETE 或业务动作）
- 中间件策略：
  - 默认拒绝。
  - 白名单仅包含登录、健康检查等公开接口。
  - 鉴权失败返回统一错误码与可观测日志。

### 1.3 完成标准
- 关键接口全部经过 Casbin 检查。
- 越权请求能被稳定拦截并记录日志。

## 1.4 第一批验收用例（后端主导，前端参与）
- 用例最少包含：
  - 正常授权访问通过
  - 未授权访问拒绝
  - 角色切换后权限立即生效
  - 订单状态变更接口的越权拦截

### 1.4 完成标准
- 自动化或半自动化用例通过率 100%。

---

## 本周建议里程碑
- D1：前端给出角色与权限矩阵初稿。
- D2：后端完成权限点映射和数据结构草案。
- D3：完成 Casbin 中间件接入与白名单配置。
- D4：跑通首批鉴权测试并修复问题。
- D5：冻结阶段 1 输出，进入阶段 2 扩展接入。

## 立即开始的第一步（今天）
- 前端输出 `角色-权限矩阵` 初稿（必须包含页面动作与 API 对应关系）。
- 后端同步输出当前 API 列表，便于前端映射。
- 双方评审后冻结第一版权限点命名规范。

## 当前进展（已启动）
- 已完成：后端 Casbin 中间件骨架接入（JWT 路由组）。
- 已完成：默认拒绝策略入口 + 白名单接口（登录、健康检查、OpenAPI）。
- 已完成：开发登录 token 注入 `sub/role`，便于前端权限方案对接前先联调。
- 待完成：前端提供第一版角色与权限矩阵，后端据此补齐策略数据与校验用例。

---

## 阶段 1.5：后端角色策略管理（与前端权限页对齐）

> 目标：支持前端“权限管理页面”直接读写后端角色策略，不再依赖本地存储。

### 1.5.1 API 设计（先冻结契约）

- `GET /central/v1/rbac/role-permissions`
  - 作用：获取全部角色的权限点列表（前端用于初始化权限页面）
  - 返回示例：
    - `roles: { super_admin: [...], manager: [...], cashier: [...], kitchen: [...], waiter: [...], unknown: [] }`

- `PUT /central/v1/rbac/role-permissions/:role`
  - 作用：更新某个角色的权限点（全量覆盖）
  - 请求体：
    - `permissions: string[]`
  - 约束：
    - role 必须在已定义角色集合中
    - permissions 必须属于已定义权限点集合

- `POST /central/v1/rbac/role-permissions/reset`
  - 作用：重置角色权限到默认模板
  - 请求体：
    - `role?: string`（不传表示重置全部角色）

### 1.5.2 数据模型设计

- 最小建议：
  - `rbac_roles`（角色主数据）
    - `id`、`role_key`、`role_name`、`is_builtin`、`created_at`、`updated_at`
  - `rbac_permissions`（权限点主数据）
    - `id`、`perm_key`、`perm_name`、`module`、`created_at`、`updated_at`
  - `rbac_role_permissions`（角色-权限关系）
    - `id`、`role_key`、`perm_key`、`created_at`、`updated_at`
  - （可选）`rbac_permission_api_map`（权限点到 API 的映射）
    - `id`、`perm_key`、`path`、`method`、`created_at`

- 唯一键建议：
  - `rbac_roles.role_key` 唯一
  - `rbac_permissions.perm_key` 唯一
  - `rbac_role_permissions (role_key, perm_key)` 联合唯一

### 1.5.3 Casbin 同步策略

- 推荐策略：
  - 以“角色”为 `sub`，即 `p, manager, /central/v1/orders, GET`
  - 用户登录后 JWT 携带 `role`

- 权限更新事务流程（核心）：
  1. 校验角色与权限点合法性
  2. 写库更新 `rbac_role_permissions`（建议事务）
  3. 根据权限点映射生成 Casbin 策略集
  4. 替换该角色的旧策略并写入新策略（同事务或补偿机制）
  5. 触发 `Enforcer.LoadPolicy()` 或增量更新

- 一致性要求：
  - 任意角色权限更新后，30 秒内全实例可见（单实例可要求立即生效）

### 1.5.4 权限点白名单（后端校验源）

- 建议由后端维护一份“可接受权限点”常量，与前端权限点同名：
  - `permission:view`
  - `home:view`
  - `workbench:view`
  - `workbench:complete`
  - `orders:view`
  - `orders:create`
  - `orders:update_status`
  - `order_desk:view`
  - `order_desk:create`
  - `menu:view`
  - `menu:edit`
  - `table:view`
  - `table:edit`

### 1.5.5 实施顺序（建议）

- D1：新增 API types + handler + logic（先返回内存假数据）
- D2：接入 DB 表与 CRUD，打通 `GET/PUT/reset`
- D3：打通 Casbin 同步与缓存刷新
- D4：联调前端权限页（确认从 local fallback 切到 remote）
- D5：补充测试与回归

### 1.5.6 验收标准

- 前端权限页首次加载显示“当前使用后端权限策略”。
- 修改任一角色权限后，刷新页面仍生效（跨浏览器/机器可见）。
- 越权接口访问能立即被 Casbin 拒绝。
- `reset` 能稳定恢复默认权限模板。
