// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSpecItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取规格项
func NewGetSpecItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpecItemLogic {
	return &GetSpecItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSpecItemLogic) GetSpecItem(req *types.GetSpecItemReq) (resp *types.GetSpecItemReply, err error) {
	item, err := l.svcCtx.SpecItem.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("规格项不存在")
		}
		return nil, err
	}

	return &types.GetSpecItemReply{
		Item: toSpecItem(item),
	}, nil
}
