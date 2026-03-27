package svc

import (
	"errors"
	"sort"
	"sync"
)

var validPermissions = map[string]struct{}{
	"permission:view":     {},
	"home:view":           {},
	"workbench:view":      {},
	"workbench:complete":  {},
	"orders:view":         {},
	"orders:create":       {},
	"orders:update_status": {},
	"order_desk:view":     {},
	"order_desk:create":   {},
	"menu:view":           {},
	"menu:edit":           {},
	"table:view":          {},
	"table:edit":          {},
}

type RbacPolicyRule struct {
	Obj string
	Act string
}

var permissionRules = map[string][]RbacPolicyRule{
	"permission:view": {
		{Obj: "/central/v1/rbac/role-permissions", Act: "GET"},
		{Obj: "/central/v1/rbac/role-permissions/super_admin", Act: "PUT"},
		{Obj: "/central/v1/rbac/role-permissions/manager", Act: "PUT"},
		{Obj: "/central/v1/rbac/role-permissions/cashier", Act: "PUT"},
		{Obj: "/central/v1/rbac/role-permissions/kitchen", Act: "PUT"},
		{Obj: "/central/v1/rbac/role-permissions/waiter", Act: "PUT"},
		{Obj: "/central/v1/rbac/role-permissions/unknown", Act: "PUT"},
		{Obj: "/central/v1/rbac/role-permissions/reset", Act: "POST"},
	},
	"workbench:view":      {{Obj: "/central/v1/workbench/orders", Act: "GET"}},
	"workbench:complete":  {{Obj: "/central/v1/order/:id/status", Act: "PUT"}},
	"orders:view":         {{Obj: "/central/v1/orders", Act: "GET"}, {Obj: "/central/v1/order/:id", Act: "GET"}},
	"orders:create":       {{Obj: "/central/v1/orders", Act: "POST"}},
	"orders:update_status": {{Obj: "/central/v1/order/:id/status", Act: "PUT"}},
	"order_desk:view": {
		{Obj: "/central/v1/menu/categories", Act: "GET"},
		{Obj: "/central/v1/menus", Act: "GET"},
		{Obj: "/central/v1/tables", Act: "GET"},
	},
	"order_desk:create": {{Obj: "/central/v1/orders", Act: "POST"}},
	"menu:view": {
		{Obj: "/central/v1/menu/categories", Act: "GET"},
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
}

type RbacStore struct {
	mu    sync.RWMutex
	roles map[string][]string
}

func NewRbacStore() *RbacStore {
	return &RbacStore{roles: defaultRolePermissions()}
}

func (s *RbacStore) List() map[string][]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneRolePermissions(s.roles)
}

func (s *RbacStore) UpdateRole(role string, permissions []string) error {
	if _, ok := s.roles[role]; !ok {
		return errors.New("role not found")
	}
	for _, permission := range permissions {
		if _, ok := validPermissions[permission]; !ok {
			return errors.New("invalid permission: " + permission)
		}
	}
	dedup := make(map[string]struct{}, len(permissions))
	for _, permission := range permissions {
		dedup[permission] = struct{}{}
	}
	next := make([]string, 0, len(dedup))
	for permission := range dedup {
		next = append(next, permission)
	}
	sort.Strings(next)

	s.mu.Lock()
	s.roles[role] = next
	s.mu.Unlock()
	return nil
}

func (s *RbacStore) Reset(role string) error {
	defaults := defaultRolePermissions()
	s.mu.Lock()
	defer s.mu.Unlock()
	if role == "" {
		s.roles = defaults
		return nil
	}
	if _, ok := s.roles[role]; !ok {
		return errors.New("role not found")
	}
	s.roles[role] = defaults[role]
	return nil
}

func BuildPoliciesForPermissions(permissions []string) []RbacPolicyRule {
	out := make([]RbacPolicyRule, 0, len(permissions)*2)
	uniq := make(map[string]struct{})
	for _, permission := range permissions {
		rules := permissionRules[permission]
		for _, rule := range rules {
			key := rule.Act + "|" + rule.Obj
			if _, ok := uniq[key]; ok {
				continue
			}
			uniq[key] = struct{}{}
			out = append(out, rule)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Obj == out[j].Obj {
			return out[i].Act < out[j].Act
		}
		return out[i].Obj < out[j].Obj
	})
	return out
}

func defaultRolePermissions() map[string][]string {
	all := make([]string, 0, len(validPermissions))
	for permission := range validPermissions {
		all = append(all, permission)
	}
	sort.Strings(all)
	return map[string][]string{
		"super_admin": append([]string{}, all...),
		"manager":     append([]string{}, all...),
		"cashier":     {"home:view", "orders:view", "orders:create", "order_desk:view", "order_desk:create"},
		"kitchen":     {"home:view", "workbench:view", "workbench:complete", "orders:view"},
		"waiter":      {"home:view", "orders:view", "order_desk:view", "order_desk:create", "table:view"},
		"unknown":     {},
	}
}

func cloneRolePermissions(in map[string][]string) map[string][]string {
	out := make(map[string][]string, len(in))
	for role, permissions := range in {
		out[role] = append([]string{}, permissions...)
	}
	return out
}
