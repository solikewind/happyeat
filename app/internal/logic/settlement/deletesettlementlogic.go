// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package settlement

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSettlementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSettlementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSettlementLogic {
	return &DeleteSettlementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSettlementLogic) DeleteSettlement(req *types.DeleteSettlementReq) (*types.DeleteSettlementReply, error) {
	if err := l.svcCtx.Settlement.Delete(l.ctx, req.Id); err != nil {
		return nil, err
	}
	return &types.DeleteSettlementReply{}, nil
}
