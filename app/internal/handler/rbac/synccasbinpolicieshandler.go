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

// 将 IAM 同步到 Casbin（刷新 casbin_rule，供管理端按钮触发）
func SyncCasbinPoliciesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SyncCasbinPoliciesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := rbac.NewSyncCasbinPoliciesLogic(r.Context(), svcCtx)
		resp, err := l.SyncCasbinPolicies(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
