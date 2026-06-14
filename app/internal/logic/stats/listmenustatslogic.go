package stats

import (
	"context"
	"sort"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	ordermodel "github.com/solikewind/happyeat/dal/model/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenuStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMenuStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuStatsLogic {
	return &ListMenuStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMenuStatsLogic) ListMenuStats(req *types.ListMenuStatsReq) (*types.ListMenuStatsReply, error) {
	r, err := ordermodel.ParseStatsDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	rows, err := l.svcCtx.Order.AggregateMenuSales(l.ctx, r)
	if err != nil {
		return nil, err
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Quantity != rows[j].Quantity {
			return rows[i].Quantity > rows[j].Quantity
		}
		return rows[i].Amount > rows[j].Amount
	})
	return buildMenuStatsReply(rows), nil
}
