package order

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/zeromicro/go-zero/core/logx"
)

type PrintOrderKitchenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPrintOrderKitchenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrintOrderKitchenLogic {
	return &PrintOrderKitchenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PrintOrderKitchenLogic) PrintOrderKitchen(req *types.PrintOrderKitchenReq) (*types.PrintOrderKitchenReply, error) {
	entOrder, err := l.svcCtx.Order.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}
	if err := SyncPrintKitchen(l.ctx, l.svcCtx, entOrder, "[手动打印]"); err != nil {
		return nil, err
	}
	return &types.PrintOrderKitchenReply{}, nil
}
