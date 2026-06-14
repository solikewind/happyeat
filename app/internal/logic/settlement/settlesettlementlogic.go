// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package settlement

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	dalsettlement "github.com/solikewind/happyeat/dal/model/settlement"

	"github.com/zeromicro/go-zero/core/logx"
)

type SettleSettlementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSettleSettlementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SettleSettlementLogic {
	return &SettleSettlementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SettleSettlementLogic) SettleSettlement(req *types.SettleSettlementReq) (*types.SettleSettlementReply, error) {
	if req.ActualAmount < 0 {
		return nil, errors.New("实收金额不能为负")
	}
	entSt, err := l.svcCtx.Settlement.Settle(l.ctx, req.Id, dalsettlement.SettleInput{
		ActualAmount: req.ActualAmount,
		Remark:       req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &types.SettleSettlementReply{
		Settlement: EntToType(l.ctx, l.svcCtx, entSt, true),
	}, nil
}
