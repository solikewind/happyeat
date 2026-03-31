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

// 获取单个角色的权限列表
func GetRolePermissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRolePermissionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := rbac.NewGetRolePermissionLogic(r.Context(), svcCtx)
		resp, err := l.GetRolePermission(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
