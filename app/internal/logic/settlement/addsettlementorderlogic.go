// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package settlement

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSettlementOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSettlementOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSettlementOrderLogic {
	return &AddSettlementOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSettlementOrderLogic) AddSettlementOrder(req *types.AddSettlementOrderReq) (*types.AddSettlementOrderReply, error) {
	if req.OrderId == 0 {
		return nil, errors.New("order_id 不能为空")
	}
	entSt, err := l.svcCtx.Settlement.AddOrder(l.ctx, req.Id, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &types.AddSettlementOrderReply{
		Settlement: EntToType(l.ctx, l.svcCtx, entSt, true),
	}, nil
}
