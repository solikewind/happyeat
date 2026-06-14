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

// 获取结账单详情（含关联订单）
func GetSettlementHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSettlementReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := settlement.NewGetSettlementLogic(r.Context(), svcCtx)
		resp, err := l.GetSettlement(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
