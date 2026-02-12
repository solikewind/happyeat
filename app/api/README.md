# HappyEat API

## 生成

```bash
make api   # 基于 central.api 一次生成所有 types、handler、logic
```

新增业务时：在 `v1/` 写 `*types.api`（可 import `base/*.api`），在 `central.api` 里 import 并加对应 `@server` 块。

## 目录结构

```
api/
├── base/                 # 实体类型（被 types 引用）
│   ├── menubase.api      # Menu, MenuCategory, MenuSpec
│   ├── tablebase.api     # Table, TableCategory
│   └── orderbase.api     # Order, OrderItem
├── v1/
│   ├── central.api       # 统一入口：import 各 types + 所有 @server 路由
│   ├── menutypes.api     # 菜单 req/reply
│   ├── tabletypes.api    # 餐桌 req/reply
│   ├── ordertypes.api    # 订单 req/reply
│   └── workbenchtypes.api # 工作台 req/reply
└── README.md
```

## 模块与路由摘要

| 模块 | 前缀/组 | 说明 |
|------|---------|------|
| menu | /menus, /menu/:id | 菜单 CRUD |
| menutype | /menu/categories, /menu/category/:id | 菜单分类 CRUD |
| table | /tables, /table/:id | 餐桌 CRUD |
| tablecategory | /table/categories, /table/category/:id | 餐桌分类 CRUD |
| order | /orders, /order/:id, /order/:id/status | 订单列表/详情/创建/更新状态 |
| workbench | /workbench/orders | 工作台订单列表（默认 created/paid/preparing）；出单用 更新订单状态 置为 completed |

## 后续扩展

- **打印订单菜品**：可增加 `GET /order/:id` 的打印用封装或 `POST /order/:id/print`，内部复用订单详情数据。
