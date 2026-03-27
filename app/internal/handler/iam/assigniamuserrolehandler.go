// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/iam"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 为用户绑定角色（幂等：已绑定则成功）；path 固定便于 Casbin obj 精确匹配
func AssignIAMUserRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AssignIAMUserRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := iam.NewAssignIAMUserRoleLogic(r.Context(), svcCtx)
		resp, err := l.AssignIAMUserRole(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
