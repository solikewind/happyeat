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

func (l *UpdateMenuCategoryLogic) UpdateMenuCategory(req *types.UpdateMenuCategoryReq) (resp *types.UpdateMenuCategoryReply, err error) {
	_, err = l.svcCtx.MenuType.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("菜单分类不存在")
		}
		return nil, err
	}
	if req.Name == "" {
		return nil, errors.New("分类名称不能为空")
	}
	err = l.svcCtx.MenuType.Update(l.ctx, req.Id, req.Name, req.Description, req.Sort)
	if err != nil {
		return nil, err
	}
	return &types.UpdateMenuCategoryReply{}, nil
}
