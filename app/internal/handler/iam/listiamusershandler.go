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

// 分页列出用户及其角色
func ListIAMUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListIAMUsersReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := iam.NewListIAMUsersLogic(r.Context(), svcCtx)
		resp, err := l.ListIAMUsers(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
