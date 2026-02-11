// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package table

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/table"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除餐桌
func DeleteTableHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteTableReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := table.NewDeleteTableLogic(r.Context(), svcCtx)
		resp, err := l.DeleteTable(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
