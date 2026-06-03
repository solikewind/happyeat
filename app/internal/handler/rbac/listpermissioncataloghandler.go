// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rbac

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/rbac"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 权限目录（权限码及对应 HTTP 接口，供前端勾选）
func ListPermissionCatalogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := rbac.NewListPermissionCatalogLogic(r.Context(), svcCtx)
		resp, err := l.ListPermissionCatalog()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
