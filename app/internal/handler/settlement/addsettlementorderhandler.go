// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package settlement

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/settlement"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 将订单加入结账单（任意订单状态，已取消除外）
func AddSettlementOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddSettlementOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := settlement.NewAddSettlementOrderLogic(r.Context(), svcCtx)
		resp, err := l.AddSettlementOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
