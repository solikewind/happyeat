// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetObjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取对象详情
func NewGetObjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetObjectLogic {
	return &GetObjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetObjectLogic) GetObject(req *types.GetObjectReq) (*types.GetObjectReply, error) {
	item, err := l.svcCtx.Object.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("对象不存在")
		}
		return nil, err
	}
	out := entObjectToType(item)
	out.Url = signedURLOrRaw(l.ctx, l.svcCtx, item.Key, out.Url)
	return &types.GetObjectReply{Object: out}, nil
}
