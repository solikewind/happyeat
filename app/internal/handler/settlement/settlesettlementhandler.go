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

// 结账单结账（未结账 → 已结账）
func SettleSettlementHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SettleSettlementReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := settlement.NewSettleSettlementLogic(r.Context(), svcCtx)
		resp, err := l.SettleSettlement(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
