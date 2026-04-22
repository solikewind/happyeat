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

type DeleteSpecItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除规格项
func NewDeleteSpecItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSpecItemLogic {
	return &DeleteSpecItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSpecItemLogic) DeleteSpecItem(req *types.DeleteSpecItemReq) (resp *types.DeleteSpecItemReply, err error) {
	// 检查是否被分类规格模板引用
	catCount, err := l.svcCtx.SpecItem.CountCategorySpecsByItem(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if catCount > 0 {
		return nil, errors.New("该规格项已被分类规格模板引用，无法删除")
	}

	// 检查是否被菜品规格直接引用
	menuCount, err := l.svcCtx.SpecItem.CountMenuSpecsByItem(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if menuCount > 0 {
		return nil, errors.New("该规格项已被菜品引用，无法删除")
	}

	if err = l.svcCtx.SpecItem.Delete(l.ctx, req.Id); err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("规格项不存在")
		}
		return nil, err
	}

	return &types.DeleteSpecItemReply{}, nil
}
