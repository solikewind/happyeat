package svc

func SyncRolePoliciesToCasbin(store *RbacStore, ce *CasbinEnforcer) error {
	enforcer := ce.Enforcer
	roles, err := store.List()
	if err != nil {
		return err
	}
	// Casbin v2.135+：RemoveFilteredPolicy 必须带 fieldValues，不能再用 (0) 表示「删全部」。
	policies, err := enforcer.GetPolicy()
	if err != nil {
		return err
	}
	if len(policies) > 0 {
		if _, err = enforcer.RemovePolicies(policies); err != nil {
			return err
		}
	}
	for role, permissions := range roles {
		policies := BuildPoliciesForPermissions(permissions)
		for _, policy := range policies {
			if _, err = enforcer.AddPolicy(role, policy.Obj, policy.Act); err != nil {
				return err
			}
		}
	}
	users, err := store.ListUserRoles()
	if err != nil {
		return err
	}
	grouping, err := enforcer.GetGroupingPolicy()
	if err != nil {
		return err
	}
	if len(grouping) > 0 {
		if _, err = enforcer.RemoveGroupingPolicies(grouping); err != nil {
			return err
		}
	}
	for userID, userRoles := range users {
		for _, role := range userRoles {
			if _, err = enforcer.AddGroupingPolicy(userID, role); err != nil {
				return err
			}
		}
	}
	return nil
}
