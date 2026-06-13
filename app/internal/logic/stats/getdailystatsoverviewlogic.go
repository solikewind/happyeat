package stats

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	ordermodel "github.com/solikewind/happyeat/dal/model/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDailyStatsOverviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDailyStatsOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDailyStatsOverviewLogic {
	return &GetDailyStatsOverviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDailyStatsOverviewLogic) GetDailyStatsOverview() (*types.ListDailyStatsReply, error) {
	r, err := ordermodel.ParseStatsDateRange("", "")
	if err != nil {
		return nil, err
	}
	daily, err := l.svcCtx.Order.AggregateDaily(l.ctx, r)
	if err != nil {
		return nil, err
	}
	summary, err := l.svcCtx.Order.AggregateSummary(l.ctx, r)
	if err != nil {
		return nil, err
	}
	return buildDailyStatsReply(daily, summary), nil
}
