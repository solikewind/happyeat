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

// 分页列出权限点（iam_permissions）
func ListIAMPermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListIAMPermissionsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := iam.NewListIAMPermissionsLogic(r.Context(), svcCtx)
		resp, err := l.ListIAMPermissions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
