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

// 按日期区间查询菜品销量明细（聚合多日后按菜品+规格汇总）
func ListMenuStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListMenuStatsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := stats.NewListMenuStatsLogic(r.Context(), svcCtx)
		resp, err := l.ListMenuStats(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
