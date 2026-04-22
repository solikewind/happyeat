// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menucategory

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除菜单种类
func NewDeleteMenuCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuCategoryLogic {
	return &DeleteMenuCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuCategoryLogic) DeleteMenuCategory(req *types.DeleteMenuCategoryReq) (resp *types.DeleteMenuCategoryReply, err error) {
	err = l.svcCtx.MenuType.Delete(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("菜单分类不存在")
		}
		return nil, err
	}
	return &types.DeleteMenuCategoryReply{}, nil
}
