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

type GetSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取规格组
func NewGetSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpecGroupLogic {
	return &GetSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSpecGroupLogic) GetSpecGroup(req *types.GetSpecGroupReq) (resp *types.GetSpecGroupReply, err error) {
	item, err := l.svcCtx.SpecGroup.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("规格组不存在")
		}
		return nil, err
	}

	return &types.GetSpecGroupReply{
		Group: toSpecGroup(item),
	}, nil
}
