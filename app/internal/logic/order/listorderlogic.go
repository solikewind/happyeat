// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"

	"github.com/solikewind/happyeat/common/consts/enum"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出订单
func NewListOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrderLogic {
	return &ListOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOrderLogic) ListOrder(req *types.ListOrderReq) (resp *types.ListOrderReply, err error) {
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	current := int(req.Current)
	if current <= 0 {
		current = 1
	}
	offset := (current - 1) * pageSize

	filter := order.ListOrdersFilter{
		Status:    enum.OrderStatus(req.Status),
		OrderType: enum.OrderType(req.OrderType),
		Offset:    offset,
		Limit:     pageSize,
	}
	if req.TableId > 0 {
		tid := req.TableId
		filter.TableID = &tid
	}

	list, total, err := l.svcCtx.Order.List(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	orders := make([]types.Order, 0, len(list))
	for _, e := range list {
		orders = append(orders, EntOrderToType(e))
	}

	return &types.ListOrderReply{Orders: orders, Total: total}, nil
}
