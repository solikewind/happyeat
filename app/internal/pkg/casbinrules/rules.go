// Package casbinrules 定义权限码与 (obj,act) 映射，供 svc 与 routecheck 共用。
package casbinrules

// PolicyRule 单条 Casbin 风格策略（角色投影前的资源点）。
type PolicyRule struct {
	Obj string
	Act string
}

// ValidPermissions 允许写入 RBAC 的权限码集合。
var ValidPermissions = map[string]struct{}{
	"permission:view":      {},
	"home:view":            {},
	"workbench:view":       {},
	"workbench:complete":   {},
	"orders:view":          {},
	"orders:create":        {},
	"orders:update_status": {},
	"order_desk:view":      {},
	"order_desk:create":    {},
	"menu:view":            {},
	"menu:edit":            {},
	"table:view":           {},
	"table:edit":           {},
	"spec:view":            {},
	"spec:edit":            {},
}

// PermissionCatalog IAM 权限点种子（描述入库）。
type PermissionSpec struct {
	Code        string
	Description string
}

var PermissionCatalog = []PermissionSpec{
	{Code: "permission:view", Description: "查看和维护角色权限配置"},
	{Code: "home:view", Description: "查看首页"},
	{Code: "workbench:view", Description: "查看工作台订单"},
	{Code: "workbench:complete", Description: "完成工作台订单"},
	{Code: "orders:view", Description: "查看订单"},
	{Code: "orders:create", Description: "创建订单"},
	{Code: "orders:update_status", Description: "更新订单状态"},
	{Code: "order_desk:view", Description: "点餐台查看"},
	{Code: "order_desk:create", Description: "点餐台创建订单"},
	{Code: "menu:view", Description: "查看菜单"},
	{Code: "menu:edit", Description: "编辑菜单"},
	{Code: "table:view", Description: "查看餐桌"},
	{Code: "table:edit", Description: "编辑餐桌"},
	{Code: "spec:view", Description: "查看规格模板"},
	{Code: "spec:edit", Description: "编辑规格模板"},
}

// PermissionRules 权限码 -> HTTP 资源点（与 Casbin 中间件 obj/act 一致）。
var PermissionRules = map[string][]PolicyRule{
	"permission:view": {
		{Obj: "/central/v1/rbac/role-permissions", Act: "GET"},
		{Obj: "/central/v1/rbac/role-permissions/:id", Act: "PUT"},
		{Obj: "/central/v1/rbac/role-permissions/reset", Act: "POST"},
		{Obj: "/central/v1/rbac/casbin/sync", Act: "POST"},
		{Obj: "/central/v1/iam/permissions", Act: "GET"},
		{Obj: "/central/v1/iam/roles", Act: "GET"},
		{Obj: "/central/v1/iam/users", Act: "GET"},
		{Obj: "/central/v1/iam/user-roles", Act: "POST"},
		{Obj: "/central/v1/iam/user-roles", Act: "DELETE"},
	},
	"workbench:view":       {{Obj: "/central/v1/workbench/orders", Act: "GET"}},
	"workbench:complete":   {{Obj: "/central/v1/order/:id/status", Act: "PUT"}},
	"orders:view":          {{Obj: "/central/v1/orders", Act: "GET"}, {Obj: "/central/v1/order/:id", Act: "GET"}},
	"orders:create":        {{Obj: "/central/v1/orders", Act: "POST"}},
	"orders:update_status": {{Obj: "/central/v1/order/:id/status", Act: "PUT"}},
	"order_desk:view": {
		{Obj: "/central/v1/menu/categories", Act: "GET"},
		{Obj: "/central/v1/menus", Act: "GET"},
		{Obj: "/central/v1/tables", Act: "GET"},
	},
	"order_desk:create": {{Obj: "/central/v1/orders", Act: "POST"}},
	"menu:view": {
		{Obj: "/central/v1/menu/categories", Act: "GET"},
		{Obj: "/central/v1/menu/category/:id", Act: "GET"},
		{Obj: "/central/v1/menus", Act: "GET"},
		{Obj: "/central/v1/menu/:id", Act: "GET"},
	},
	"menu:edit": {
		{Obj: "/central/v1/menu/category", Act: "POST"},
		{Obj: "/central/v1/menu/category/:id", Act: "PUT"},
		{Obj: "/central/v1/menu/category/:id", Act: "DELETE"},
		{Obj: "/central/v1/menus", Act: "POST"},
		{Obj: "/central/v1/menu/:id", Act: "PUT"},
		{Obj: "/central/v1/menu/:id", Act: "DELETE"},
	},
	"table:view": {
		{Obj: "/central/v1/table/categories", Act: "GET"},
		{Obj: "/central/v1/table/category/:id", Act: "GET"},
		{Obj: "/central/v1/tables", Act: "GET"},
		{Obj: "/central/v1/table/:id", Act: "GET"},
	},
	"table:edit": {
		{Obj: "/central/v1/table/category", Act: "POST"},
		{Obj: "/central/v1/table/category/:id", Act: "PUT"},
		{Obj: "/central/v1/table/category/:id", Act: "DELETE"},
		{Obj: "/central/v1/tables", Act: "POST"},
		{Obj: "/central/v1/table/:id", Act: "PUT"},
		{Obj: "/central/v1/table/:id", Act: "DELETE"},
	},
	"spec:view": {
		{Obj: "/central/v1/spec/category-spec", Act: "GET"},
		{Obj: "/central/v1/spec/category-spec/:id", Act: "GET"},
		{Obj: "/central/v1/spec/group/:id", Act: "GET"},
		{Obj: "/central/v1/spec/groups", Act: "GET"},
		{Obj: "/central/v1/spec/item/:id", Act: "GET"},
		{Obj: "/central/v1/spec/items", Act: "GET"},
	},
	"spec:edit": {
		{Obj: "/central/v1/spec/category-spec", Act: "POST"},
		{Obj: "/central/v1/spec/category-spec/:id", Act: "PUT"},
		{Obj: "/central/v1/spec/category-spec/:id", Act: "DELETE"},
		{Obj: "/central/v1/spec/group", Act: "POST"},
		{Obj: "/central/v1/spec/group/:id", Act: "PUT"},
		{Obj: "/central/v1/spec/group/:id", Act: "DELETE"},
		{Obj: "/central/v1/spec/item", Act: "POST"},
		{Obj: "/central/v1/spec/item/:id", Act: "PUT"},
		{Obj: "/central/v1/spec/item/:id", Act: "DELETE"},
	},
}

// AllPolicyKeys 返回所有 (canonicalObj, act) 用于与路由集合 diff。
func AllPolicyKeys(canonical func(obj, act string) string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, rules := range PermissionRules {
		for _, r := range rules {
			k := canonical(r.Obj, r.Act)
			out[k] = struct{}{}
		}
	}
	return out
}
