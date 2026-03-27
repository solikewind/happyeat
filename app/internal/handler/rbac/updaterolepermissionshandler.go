// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rbac

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/rbac"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新角色权限（全量覆盖）
func UpdateRolePermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateRolePermissionsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := rbac.NewUpdateRolePermissionsLogic(r.Context(), svcCtx)
		resp, err := l.UpdateRolePermissions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
