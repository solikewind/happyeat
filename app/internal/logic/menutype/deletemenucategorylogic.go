// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menutype

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/menu/ent"

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

func (l *DeleteMenuCategoryLogic) DeleteMenuCategory(req *types.DeleteMenuCategoryReq) (*types.DeleteMenuCategoryReply, error) {
	count, err := l.svcCtx.MenuType.CountMenusByCategoryID(l.ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该分类下仍有菜单，无法删除")
	}

	err = l.svcCtx.MenuType.Delete(l.ctx, int(req.Id))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, err
		}
		return nil, err
	}

	return &types.DeleteMenuCategoryReply{}, nil
}
