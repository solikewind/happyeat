// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menutype

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新菜单种类
func NewUpdateMenuCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuCategoryLogic {
	return &UpdateMenuCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuCategoryLogic) UpdateMenuCategory(req *types.UpdateMenuCategoryReq) (*types.UpdateMenuCategoryReply, error) {
	c := req.MenuCategory

	err := l.svcCtx.MenuType.Update(l.ctx, int(req.Id), c.Name, c.Description)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, err
		}
		return nil, err
	}

	return &types.UpdateMenuCategoryReply{}, nil
}
