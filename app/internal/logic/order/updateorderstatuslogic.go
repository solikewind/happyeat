// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/constants"
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
var allowedTransitions = map[string]map[string]bool{
	constants.OrderStatusCreated: {
		constants.OrderStatusPaid: true, constants.OrderStatusCancelled: true,
	},
	constants.OrderStatusPaid: {
		constants.OrderStatusPreparing: true, constants.OrderStatusCancelled: true,
	},
	constants.OrderStatusPreparing: {
		constants.OrderStatusCompleted: true,
	},
	// completed / cancelled 为终态，不可再改
}

func (l *UpdateOrderStatusLogic) UpdateOrderStatus(req *types.UpdateOrderStatusReq) (*types.UpdateOrderStatusReply, error) {
	cur, err := l.svcCtx.Order.GetByID(l.ctx, int(req.Id))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}

	next := req.Status
	allowed, ok := allowedTransitions[cur.Status]
	if !ok || !allowed[next] {
		return nil, errors.New("当前状态不允许变更为 " + next)
	}

	err = l.svcCtx.Order.UpdateStatus(l.ctx, int(req.Id), req.Status)
	if err != nil {
		return nil, err
	}

	return &types.UpdateOrderStatusReply{}, nil
}
