// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package settlement

import (
	"context"
	"errors"
	"strings"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSettlementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSettlementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSettlementLogic {
	return &CreateSettlementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSettlementLogic) CreateSettlement(req *types.CreateSettlementReq) (*types.CreateSettlementReply, error) {
	if strings.TrimSpace(req.CustomerName) == "" {
		return nil, errors.New("客户名不能为空")
	}
	entSt, err := l.svcCtx.Settlement.Create(l.ctx, req.CustomerName, req.Remark)
	if err != nil {
		return nil, err
	}
	return &types.CreateSettlementReply{
		Settlement: EntToType(l.ctx, l.svcCtx, entSt, false),
	}, nil
}
