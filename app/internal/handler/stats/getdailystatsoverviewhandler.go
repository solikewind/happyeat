// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package stats

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/stats"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 今日经营概览（等价于 start_date=end_date=今天）
func GetDailyStatsOverviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := stats.NewGetDailyStatsOverviewLogic(r.Context(), svcCtx)
		resp, err := l.GetDailyStatsOverview()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
