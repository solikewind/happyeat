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

// 移除用户的某个角色；使用 query：?user_code=&role_code=
func RemoveIAMUserRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveIAMUserRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := iam.NewRemoveIAMUserRoleLogic(r.Context(), svcCtx)
		resp, err := l.RemoveIAMUserRole(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
