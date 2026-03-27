package svc

func SyncRolePoliciesToCasbin(store *RbacStore, ce *CasbinEnforcer) error {
	enforcer := ce.Enforcer
	roles := store.List()
	for role, permissions := range roles {
		policies := BuildPoliciesForPermissions(permissions)
		if _, err := enforcer.RemoveFilteredPolicy(0, role); err != nil {
			return err
		}
		for _, policy := range policies {
			if _, err := enforcer.AddPolicy(role, policy.Obj, policy.Act); err != nil {
				return err
			}
		}
	}
	return nil
}
