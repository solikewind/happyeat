# Casbin RBAC 设计说明（单租户）

本文档说明 `happyeat` 当前权限系统的最小实现与使用方式。

## 1. 为什么用 Casbin

Casbin 负责运行时鉴权判定，核心优势是：

- 策略结构统一（`p` / `g`）。
- 接口鉴权链路固定（`Enforce(sub, obj, act)`）。
- 适合后续扩展为多租户（domain）模型。

## 2. 当前采用的模型

当前使用单租户 RBAC：

- request: `r = sub, obj, act`
- policy: `p = sub, obj, act`
- role: `g = _, _`
- matcher: `g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`

说明：

- `sub`：用户 ID（来自 JWT 的 `sub`）。
- `obj`：标准化后的接口路径（如 `/central/v1/order/:id`）。
- `act`：HTTP 方法大写（`GET/POST/PUT/DELETE`）。

### 与 vega「operation code」模式的差异

- **vega**：`obj` 通常对应 **operation code**（业务操作码），与具体 HTTP 路由解耦。
- **本项目（go-zero REST）**：`obj` 采用 **归一化后的 HTTP 路径**（与中间件里 `Enforce` 使用的字符串一致）。主存权限码（如 `menu:view`）到 Casbin 的投影思路与 vega 一致，只是 **资源标识**按栈选定为 REST path。

路径归一化与「路由模板 ↔ 权限表」对齐逻辑在 `app/internal/pkg/routenorm`；权限码到 `(obj, act)` 的映射集中在 `app/internal/pkg/casbinrules`，由 `rbac_store` 与校验工具共用，避免双处维护漂移。

## 3. 主存与投影

权限数据分两层：

1) 业务主存（真相源）

- `iam_users`
- `iam_roles`
- `iam_permissions`
- `iam_user_roles`
- `iam_role_permissions`

2) Casbin 投影（运行时）

- `casbin_rule`
  - `p`：角色对资源动作的授权
  - `g`：用户与角色的继承关系

约束：业务侧只维护主存表，Casbin 由同步器统一生成，避免双写不一致。

## 4. 鉴权链路

1. 登录签发 JWT，仅携带 `sub`（用户 ID）及时间字段。  
2. 中间件解析 JWT，提取 `sub`。  
3. 中间件归一化路径为 `obj`，提取方法为 `act`。  
4. 调用 `Enforce(sub, obj, act)` 判定放行或拒绝。  

## 5. 权限点规范

当前 `permission` 采用业务编码（如 `orders:view`、`menu:edit`），在服务内映射到一组 `obj/act`。

建议规范：

- 命名：`{module}:{action}`，例如 `table:view`、`spec:edit`。
- 动作统一词汇：`view/create/edit/update_status/complete`。
- 新增接口必须同时补充到权限映射，保证策略可投影。

## 6. 策略同步时机

以下场景触发同步：

- 服务启动（首次装载）。
- 更新角色权限（`PUT /rbac/role-permissions/:role`）。
- 重置角色权限（`POST /rbac/role-permissions/reset`）。

同步内容：

- 先重建全部 `p` 规则（角色 -> 资源动作）。
- 再重建全部 `g` 规则（用户 -> 角色）。

## 7. 新增权限接入步骤

1. 新增权限编码（如 `orders:refund`）。  
2. 维护该编码对应的 `obj/act` 映射。  
3. 把权限分配给目标角色。  
4. 触发策略同步。  
5. 用目标账号访问接口验证 `2xx/403`。  

## 8. 发版前：`routecheck` 路由与权限映射校验

命令 `app/cmd/routecheck` 会：

- 用 `go/ast` 解析 `app/internal/handler/routes.go` 中挂了 `CasbinMiddleware` 的 `AddRoutes` 块，拼出 `WithPrefix` + `Path` 与 HTTP Method；
- 与 `casbinrules.PermissionRules` 聚合出的 `(obj, act)` 集合对比（对比键对路径参数做了与工具一致的规范化，例如 `:role` 与 `:id` 视为同类槽位）。

**在模块根目录 `happyeat` 执行：**

```bash
go run ./app/cmd/routecheck
```

- 若存在「路由已受 Casbin 保护、但未出现在任何 permission 的 `obj/act` 映射中」的项，命令以 **退出码 1** 结束（高风险：默认拒绝）。
- 「映射里有、当前解析到的 Casbin 路由里没有」的项会报告为 **僵尸规则**；默认仍 **退出码 0**，若希望在 CI 中一并失败可加 `-strict`。

其它常用参数：

- `-json`：输出 JSON，便于流水线解析。
- `-routes <path>`：指定非默认的 `routes.go` 路径。

发版或合并前建议执行 `make routecheck`（或上述 `go run`），并与 `go test ./...` 一起作为门禁。

## 9. 向多租户演进（后续）

后续可平滑升级为 domain 模型：

- 模型升级为 `r=sub, dom, obj, act` / `p=sub, dom, obj, act` / `g=_, _, _`。
- 在主存表增加 `tenant_id` 维度。
- JWT 增加租户字段，中间件把租户映射到 `dom`。

当前单租户实现已将“主存与投影”职责拆分，升级时只需扩展维度，不需要推翻结构。
