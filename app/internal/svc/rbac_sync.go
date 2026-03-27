package svc

func SyncRolePoliciesToCasbin(store *RbacStore, ce *CasbinEnforcer) error {
	enforcer := ce.Enforcer
	roles, err := store.List()
	if err != nil {
		return err
	}
	if _, err = enforcer.RemoveFilteredPolicy(0); err != nil {
		return err
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
	if _, err = enforcer.RemoveFilteredGroupingPolicy(0); err != nil {
		return err
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
