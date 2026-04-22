# HappyEat 缺陷与踩坑记录

本文件记录「已定位原因 + 已修复或已规避」的线上/开发问题，便于复盘与避免重复踩坑。  
**新增条目**：修复或确认某类问题后，在表尾追加一行，并写清日期、现象、根因、处理。

## 记录表

| 日期 | 模块 | 现象（摘要） | 根因 | 处理 / 备注 |
|------|------|--------------|------|-------------|
| 2026-04-19 | 订单 `Order.status` | `POST /central/v1/orders` 返回 500；日志：`ent: validator failed for field "Order.status": order: invalid enum value for status field: "created"` | **ent 枚举**在 `common/consts/enum` 中定义为 **大写**（`CREATED`、`PAID` …），`StatusValidator` 只接受大写。若某条路径传入 **小写** `"created"`（例如旧代码 `SetStatus("created")`、`pkg/status` 曾用小写常量、或未重启的旧二进制），校验在 `Save` 前即失败。 | 与 `enum` 统一为大写；状态机 `stateless` 状态名与 `ORDER_*` 常量对齐；**DAL** `dal/model/order/order.go` 增加 `normalizeOrderStatus`，在 `Create` / `UpdateStatus` 的 `SetStatus` 前将合法小写统一映射为大写枚举，避免上游偶发小写。创建订单时 `Status` 使用 `enum.OrderStatusCreated`。 |

## 相关代码位置（订单状态）

- 枚举：`common/consts/enum/order_status.go`
- ent 校验：`dal/model/ent/order/order.go`（生成代码中的 `StatusValidator`）
- 业务 DAL 封装：`dal/model/order/order.go`（`normalizeOrderStatus`、`Create`、`UpdateStatus`）
- 状态机：`app/internal/pkg/status/order_statemachine.go`、`order_transition.go`
- 更新状态 logic：`app/internal/logic/order/updateorderstatuslogic.go`
