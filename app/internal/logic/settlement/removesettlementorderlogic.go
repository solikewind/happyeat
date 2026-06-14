// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package settlement

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveSettlementOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveSettlementOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveSettlementOrderLogic {
	return &RemoveSettlementOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveSettlementOrderLogic) RemoveSettlementOrder(req *types.RemoveSettlementOrderReq) (*types.RemoveSettlementOrderReply, error) {
	entSt, err := l.svcCtx.Settlement.RemoveOrder(l.ctx, req.Id, req.OrderId)
	if err != nil {
		return nil, err
	}
	return &types.RemoveSettlementOrderReply{
		Settlement: EntToType(l.ctx, l.svcCtx, entSt, true),
	}, nil
}
