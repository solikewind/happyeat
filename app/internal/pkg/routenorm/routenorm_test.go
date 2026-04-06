package routenorm

import "testing"

func TestEnforceObj(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"/central/v1/rbac/role-permissions/super_admin", "/central/v1/rbac/role-permissions/:id"},
		{"/central/v1/rbac/role-permissions/reset", "/central/v1/rbac/role-permissions/reset"},
		{"/central/v1/rbac/role-permissions/reset/", "/central/v1/rbac/role-permissions/reset"},
		{"/central/v1/order/42/status", "/central/v1/order/:id/status"},
		{"/central/v1/order/42", "/central/v1/order/:id"},
		{"/central/v1/menu/550e8400-e29b-41d4-a716-446655440000", "/central/v1/menu/:id"},
		{"/central/v1/iam/permissions", "/central/v1/iam/permissions"},
	}
	for _, tt := range tests {
		if got := EnforceObj(tt.in); got != tt.want {
			t.Errorf("EnforceObj(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
