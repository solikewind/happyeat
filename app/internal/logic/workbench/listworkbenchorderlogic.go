// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package workbench

import (
	"context"
	"strings"

	orderlogic "github.com/solikewind/happyeat/app/internal/logic/order"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/consts/enum"
	"github.com/solikewind/happyeat/dal/model/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkbenchOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 工作台订单列表（默认待处理：created/paid/preparing）；出单用 更新订单状态 置为 completed
func NewListWorkbenchOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkbenchOrderLogic {
	return &ListWorkbenchOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWorkbenchOrderLogic) ListWorkbenchOrder(req *types.ListWorkbenchOrderReq) (resp *types.ListWorkbenchOrderReply, err error) {
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
		Offset: offset,
		Limit:  pageSize,
	}
	if strings.TrimSpace(req.Status) != "" {
		filter.Status = enum.OrderStatus(strings.TrimSpace(req.Status))
	} else {
		filter.Statuses = []enum.OrderStatus{
			enum.OrderStatusCreated,
			enum.OrderStatusPaid,
			enum.OrderStatusPreparing,
		}
	}

	list, total, err := l.svcCtx.Order.List(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	orders := make([]types.Order, 0, len(list))
	for _, e := range list {
		orders = append(orders, orderlogic.EntOrderToType(e))
	}

	return &types.ListWorkbenchOrderReply{Orders: orders, Total: total}, nil
}
