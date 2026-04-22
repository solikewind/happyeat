// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package order

import (
	"context"
	"errors"

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

func (l *UpdateOrderStatusLogic) UpdateOrderStatus(req *types.UpdateOrderStatusReq) (*types.UpdateOrderStatusReply, error) {
	cur, err := l.svcCtx.Order.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}

	curMachine := status.EntStatusToMachine(cur.Status)
	nextMachine, err := status.ParseAPIStatus(req.Status)
	if err != nil {
		return nil, err
	}

	trigger, err := status.ResolveTrigger(curMachine, nextMachine)
	if err != nil {
		return nil, err
	}

	sm := status.NewOrderStateMachine(curMachine, *cur)
	if err := sm.FireCtx(l.ctx, trigger); err != nil {
		return nil, err
	}

	newMachine, err := sm.CurrentMachineState(l.ctx)
	if err != nil {
		return nil, err
	}
	nextEnum, err := status.MachineToEntStatus(newMachine)
	if err != nil {
		return nil, err
	}

	if err := l.svcCtx.Order.UpdateStatus(l.ctx, req.Id, nextEnum); err != nil {
		return nil, err
	}

	return &types.UpdateOrderStatusReply{}, nil
}
