// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stats

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/stats"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 按日期区间查询日汇总（start_date/end_date 格式 YYYY-MM-DD；均可缺省为今天；可查近3/7/30天等）
func ListDailyStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListDailyStatsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := stats.NewListDailyStatsLogic(r.Context(), svcCtx)
		resp, err := l.ListDailyStats(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
