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

// 更新餐桌
func UpdateTableHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTableReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := table.NewUpdateTableLogic(r.Context(), svcCtx)
		resp, err := l.UpdateTable(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
