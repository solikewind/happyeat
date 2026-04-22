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

type DeleteCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除分类规格模板
func NewDeleteCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategorySpecLogic {
	return &DeleteCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCategorySpecLogic) DeleteCategorySpec(req *types.DeleteCategorySpecReq) (resp *types.DeleteCategorySpecReply, err error) {
	menuCount, err := l.svcCtx.CategorySpec.CountMenuSpecsByCategory(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if menuCount > 0 {
		return nil, errors.New("该分类规格模板已被菜品引用，无法删除")
	}

	if err = l.svcCtx.CategorySpec.Delete(l.ctx, req.Id); err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("分类规格不存在")
		}
		return nil, err
	}

	return &types.DeleteCategorySpecReply{}, nil
}
