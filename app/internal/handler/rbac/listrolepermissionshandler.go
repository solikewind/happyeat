// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rbac

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/rbac"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取角色权限矩阵
func ListRolePermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := rbac.NewListRolePermissionsLogic(r.Context(), svcCtx)
		resp, err := l.ListRolePermissions()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
