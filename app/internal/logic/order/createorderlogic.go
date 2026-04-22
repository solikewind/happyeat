// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/common/consts/enum"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建订单
func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (*types.CreateOrderReply, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("订单至少需要一项")
	}
	if req.OrderType != "dine_in" && req.OrderType != "takeaway" {
		return nil, errors.New("order_type 应为 dine_in 或 takeaway")
	}
	if req.ActualAmount < 0 {
		return nil, errors.New("实收金额不能为负")
	}
	actualAmount := req.ActualAmount
	if actualAmount == 0 && req.TotalAmount > 0 {
		actualAmount = req.TotalAmount
	}

	items := make([]order.ItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, order.ItemInput{
			MenuName:  it.MenuName,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
			SpecInfo:  it.SpecInfo,
		})
	}

	var tableID *uint64
	if req.TableId > 0 {
		tid := req.TableId
		tableID = &tid
	}

	// 初始状态：CREATED（与 pkg/status 状态机 NONE→TriggerCreate 一致；此处直接写枚举，避免请求内再 Fire 一遍）
	entOrder, err := l.svcCtx.Order.Create(l.ctx, order.CreateOrderInput{
		OrderType:    enum.OrderType(req.OrderType),
		TableID:      tableID,
		Items:        items,
		TotalAmount:  req.TotalAmount,
		ActualAmount: actualAmount,
		Remark:       req.Remark,
		Status:       enum.OrderStatusCreated,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateOrderReply{Order: EntOrderToType(entOrder)}, nil
}
