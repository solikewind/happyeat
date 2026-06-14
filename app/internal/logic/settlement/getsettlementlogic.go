// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package settlement

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSettlementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSettlementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSettlementLogic {
	return &GetSettlementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSettlementLogic) GetSettlement(req *types.GetSettlementReq) (*types.GetSettlementReply, error) {
	entSt, err := l.svcCtx.Settlement.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("结账单不存在")
		}
		return nil, err
	}
	return &types.GetSettlementReply{
		Settlement: EntToType(l.ctx, l.svcCtx, entSt, true),
	}, nil
}
