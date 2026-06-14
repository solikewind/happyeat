package stats

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	ordermodel "github.com/solikewind/happyeat/dal/model/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDailyStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDailyStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDailyStatsLogic {
	return &ListDailyStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDailyStatsLogic) ListDailyStats(req *types.ListDailyStatsReq) (*types.ListDailyStatsReply, error) {
	r, err := ordermodel.ParseStatsDateRange(req.StartDate, req.EndDate)
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
