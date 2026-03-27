package svc

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/solikewind/happyeat/app/internal/pkg/casbinrules"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/ent/iampermission"
	"github.com/solikewind/happyeat/dal/model/ent/iamrole"
	"github.com/solikewind/happyeat/dal/model/ent/iamuser"
)

// RbacPolicyRule 与 Casbin 投影中的 (obj, act) 一致，定义见 casbinrules.PolicyRule。
type RbacPolicyRule = casbinrules.PolicyRule

type RbacStore struct {
	mu     sync.RWMutex
	client *ent.Client
}

func NewRbacStore(client *ent.Client) (*RbacStore, error) {
	store := &RbacStore{client: client}
	if err := store.bootstrap(context.Background()); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *RbacStore) List() (map[string][]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ctx := context.Background()
	roles, err := s.client.IAMRole.Query().
		WithPermissions(func(q *ent.IAMPermissionQuery) {
			q.Order(ent.Asc(iampermission.FieldPermissionCode))
		}).
		Order(ent.Asc(iamrole.FieldRoleCode)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make(map[string][]string, len(roles))
	for _, role := range roles {
		out[role.RoleCode] = make([]string, 0, len(role.Edges.Permissions))
		for _, permission := range role.Edges.Permissions {
			out[role.RoleCode] = append(out[role.RoleCode], permission.PermissionCode)
		}
	}
	return out, nil
}

func (s *RbacStore) UpdateRole(roleCode string, permissions []string) error {
	ctx := context.Background()
	roleEnt, err := s.client.IAMRole.Query().Where(iamrole.RoleCodeEQ(roleCode)).Only(ctx)
	if err != nil {
		return errors.New("role not found")
	}

	seen := map[string]struct{}{}
	dedup := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		if _, ok := casbinrules.ValidPermissions[permission]; !ok {
			return errors.New("invalid permission: " + permission)
		}
		if _, ok := seen[permission]; ok {
			continue
		}
		seen[permission] = struct{}{}
		dedup = append(dedup, permission)
	}
	sort.Strings(dedup)

	tx, err := s.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updater := tx.IAMRole.UpdateOneID(roleEnt.ID).ClearPermissions()
	if len(dedup) > 0 {
		permissionEnts, err := tx.IAMPermission.Query().
			Where(iampermission.PermissionCodeIn(dedup...)).
			All(ctx)
		if err != nil {
			return err
		}
		if len(permissionEnts) != len(dedup) {
			return errors.New("permissions not found")
		}
		permissionIDs := make([]uint64, 0, len(permissionEnts))
		for _, item := range permissionEnts {
			permissionIDs = append(permissionIDs, item.ID)
		}
		updater = updater.AddPermissionIDs(permissionIDs...)
	}
	if _, err = updater.Save(ctx); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RbacStore) ListUserRoles() (map[string][]string, error) {
	ctx := context.Background()
	users, err := s.client.IAMUser.Query().
		WithRoles(func(q *ent.IAMRoleQuery) {
			q.Order(ent.Asc(iamrole.FieldRoleCode))
		}).
		Order(ent.Asc(iamuser.FieldUserCode)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make(map[string][]string, len(users))
	for _, user := range users {
		out[user.UserCode] = make([]string, 0, len(user.Edges.Roles))
		for _, role := range user.Edges.Roles {
			out[user.UserCode] = append(out[user.UserCode], role.RoleCode)
		}
	}
	return out, nil
}

func (s *RbacStore) EnsureUser(userCode string) error {
	if userCode == "" {
		return errors.New("user_code 不能为空")
	}
	ctx := context.Background()
	exists, err := s.client.IAMUser.Query().Where(iamuser.UserCodeEQ(userCode)).Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = s.client.IAMUser.Create().
		SetUserCode(userCode).
		SetDisplayName(userCode).
		Save(ctx)
	return err
}

func (s *RbacStore) AssignUserRole(userCode, roleCode string) error {
	if err := s.requireRole(roleCode); err != nil {
		return err
	}
	if err := s.EnsureUser(userCode); err != nil {
		return err
	}

	ctx := context.Background()
	userEnt, err := s.client.IAMUser.Query().Where(iamuser.UserCodeEQ(userCode)).Only(ctx)
	if err != nil {
		return err
	}
	roleEnt, err := s.client.IAMRole.Query().Where(iamrole.RoleCodeEQ(roleCode)).Only(ctx)
	if err != nil {
		return err
	}
	hasRole, err := s.client.IAMUser.Query().
		Where(iamuser.UserCodeEQ(userCode), iamuser.HasRolesWith(iamrole.RoleCodeEQ(roleCode))).
		Exist(ctx)
	if err != nil {
		return err
	}
	if hasRole {
		return nil
	}
	_, err = s.client.IAMUser.UpdateOneID(userEnt.ID).AddRoleIDs(roleEnt.ID).Save(ctx)
	return err
}

func (s *RbacStore) Reset(roleCode string) error {
	defaults := defaultRolePermissions()
	if roleCode == "" {
		for rc, permissions := range defaults {
			if err := s.UpdateRole(rc, permissions); err != nil {
				return err
			}
		}
		return nil
	}
	permissions, ok := defaults[roleCode]
	if !ok {
		return errors.New("role not found")
	}
	return s.UpdateRole(roleCode, permissions)
}

func (s *RbacStore) bootstrap(ctx context.Context) error {
	if err := s.seedRolesAndPermissions(ctx); err != nil {
		return err
	}
	if err := s.seedDefaultMappings(); err != nil {
		return err
	}
	return nil
}

func (s *RbacStore) seedRolesAndPermissions(ctx context.Context) error {
	defaults := defaultRolePermissions()
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for roleCode := range defaults {
		exists, err := tx.IAMRole.Query().Where(iamrole.RoleCodeEQ(roleCode)).Exist(ctx)
		if err != nil {
			return err
		}
		if !exists {
			if _, err = tx.IAMRole.Create().SetRoleCode(roleCode).SetRoleName(roleCode).Save(ctx); err != nil {
				return err
			}
		}
	}
	for _, permission := range casbinrules.PermissionCatalog {
		exists, err := tx.IAMPermission.Query().Where(iampermission.PermissionCodeEQ(permission.Code)).Exist(ctx)
		if err != nil {
			return err
		}
		if !exists {
			if _, err = tx.IAMPermission.Create().
				SetPermissionCode(permission.Code).
				SetDescription(permission.Description).
				Save(ctx); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *RbacStore) seedDefaultMappings() error {
	ctx := context.Background()
	count, err := s.client.IAMRole.Query().QueryPermissions().Count(ctx)
	if err != nil {
		return err
	}
	if count == 0 {
		defaults := defaultRolePermissions()
		for roleCode, permissions := range defaults {
			if err := s.UpdateRole(roleCode, permissions); err != nil {
				return err
			}
		}
	}
	if err := s.EnsureUser("dev-admin"); err != nil {
		return err
	}
	hasRole, err := s.client.IAMUser.Query().
		Where(iamuser.UserCodeEQ("dev-admin"), iamuser.HasRolesWith(iamrole.RoleCodeEQ("super_admin"))).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !hasRole {
		if err := s.AssignUserRole("dev-admin", "super_admin"); err != nil {
			return err
		}
	}
	return nil
}

func (s *RbacStore) requireRole(roleCode string) error {
	ctx := context.Background()
	exists, err := s.client.IAMRole.Query().Where(iamrole.RoleCodeEQ(roleCode)).Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("role not found")
	}
	return nil
}

func BuildPoliciesForPermissions(permissions []string) []RbacPolicyRule {
	out := make([]RbacPolicyRule, 0, len(permissions)*2)
	uniq := make(map[string]struct{})
	for _, permission := range permissions {
		rules := casbinrules.PermissionRules[permission]
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
	all := make([]string, 0, len(casbinrules.ValidPermissions))
	for permission := range casbinrules.ValidPermissions {
		all = append(all, permission)
	}
	sort.Strings(all)
	return map[string][]string{
		"super_admin": append([]string{}, all...),
		"manager":     append([]string{}, all...),
		"cashier":     {"home:view", "orders:view", "orders:create", "order_desk:view", "order_desk:create", "spec:view"},
		"kitchen":     {"home:view", "workbench:view", "workbench:complete", "orders:view", "spec:view"},
		"waiter":      {"home:view", "orders:view", "order_desk:view", "order_desk:create", "table:view", "spec:view"},
		"unknown":     {},
	}
}
