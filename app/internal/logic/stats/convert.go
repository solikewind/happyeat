package stats

import (
	ordermodel "github.com/solikewind/happyeat/dal/model/order"
	"github.com/solikewind/happyeat/app/internal/types"
)

const dateLayout = "2006-01-02"

func buildDailyStatsReply(daily []ordermodel.DailyStatsPoint, summary ordermodel.StatsSummary) *types.ListDailyStatsReply {
	points := make([]types.DailyStatsPoint, 0, len(daily))
	for _, p := range daily {
		points = append(points, types.DailyStatsPoint{
			Date:          p.Date.Format(dateLayout),
			OrderCount:    p.OrderCount,
			Revenue:       p.Revenue,
			ItemCount:     p.ItemCount,
			DineInCount:   p.DineInCount,
			TakeawayCount: p.TakeawayCount,
		})
	}
	return &types.ListDailyStatsReply{
		Daily:   points,
		Summary: toSummaryType(summary),
	}
}

func toSummaryType(s ordermodel.StatsSummary) types.DailyStatsSummary {
	return types.DailyStatsSummary{
		OrderCount:    s.OrderCount,
		Revenue:       s.Revenue,
		ItemCount:     s.ItemCount,
		DineInCount:   s.DineInCount,
		TakeawayCount: s.TakeawayCount,
	}
}

func buildMenuStatsReply(rows []ordermodel.MenuSalesRow) *types.ListMenuStatsReply {
	out := make([]types.MenuStatsRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, types.MenuStatsRow{
			MenuName: r.MenuName,
			SpecInfo: r.SpecInfo,
			Quantity: r.Quantity,
			Amount:   r.Amount,
		})
	}
	return &types.ListMenuStatsReply{Rows: out}
}
