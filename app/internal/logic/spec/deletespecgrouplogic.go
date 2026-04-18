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

type DeleteSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除规格组
func NewDeleteSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSpecGroupLogic {
	return &DeleteSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSpecGroupLogic) DeleteSpecGroup(req *types.DeleteSpecGroupReq) (resp *types.DeleteSpecGroupReply, err error) {
	itemCount, err := l.svcCtx.SpecGroup.CountItemsByGroupID(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if itemCount > 0 {
		return nil, errors.New("规格组下仍有规格项，无法删除")
	}

	if err = l.svcCtx.SpecGroup.Delete(l.ctx, req.Id); err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("规格组不存在")
		}
		return nil, err
	}

	return &types.DeleteSpecGroupReply{}, nil
}
