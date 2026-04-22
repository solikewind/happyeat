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

type GetMenuCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取菜单种类
func NewGetMenuCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuCategoryLogic {
	return &GetMenuCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuCategoryLogic) GetMenuCategory(req *types.GetMenuCategoryReq) (resp *types.GetMenuCategoryReply, err error) {
	category, err := l.svcCtx.MenuType.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("菜单分类不存在")
		}
		return nil, err
	}
	return &types.GetMenuCategoryReply{
		Category: entMenuCategoryToType(category),
	}, nil
}
