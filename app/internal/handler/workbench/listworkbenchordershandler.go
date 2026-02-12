// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package workbench

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/workbench"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 工作台订单列表（默认待处理：created/paid/preparing）；出单用 更新订单状态 置为 completed
func ListWorkbenchOrdersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListWorkbenchOrdersReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := workbench.NewListWorkbenchOrdersLogic(r.Context(), svcCtx)
		resp, err := l.ListWorkbenchOrders(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
