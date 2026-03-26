// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/common/consts/enum"
	"github.com/solikewind/happyeat/app/internal/pkg/status"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新订单状态
func NewUpdateOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 允许的状态流转：当前状态 -> 可变更到的状态
var allowedTransitions = map[enum.OrderStatus]map[enum.OrderStatus]bool{
	enum.OrderStatus(status.OrderStatusCreated): {
		enum.OrderStatus(status.OrderStatusPaid): true, enum.OrderStatus(status.OrderStatusCancelled): true,
	},
	enum.OrderStatus(status.OrderStatusPaid): {
		enum.OrderStatus(status.OrderStatusPreparing): true, enum.OrderStatus(status.OrderStatusCancelled): true,
	},
	enum.OrderStatus(status.OrderStatusPreparing): {
		enum.OrderStatus(status.OrderStatusCompleted): true,
	},
	// completed / cancelled 为终态，不可再改
}

func (l *UpdateOrderStatusLogic) UpdateOrderStatus(req *types.UpdateOrderStatusReq) (*types.UpdateOrderStatusReply, error) {
	cur, err := l.svcCtx.Order.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}

	next := enum.OrderStatus(req.Status)
	allowed, ok := allowedTransitions[cur.Status]
	if !ok || !allowed[next] {
		return nil, errors.New("当前状态不允许变更为 " + string(next))
	}

	err = l.svcCtx.Order.UpdateStatus(l.ctx, req.Id, enum.OrderStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &types.UpdateOrderStatusReply{}, nil
}
