// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menucategory

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/menu"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建菜单种类
func NewCreateMenuCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuCategoryLogic {
	return &CreateMenuCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuCategoryLogic) CreateMenuCategory(req *types.CreateMenuCategoryReq) (resp *types.CreateMenuCategoryReply, err error) {
	if req.Name == "" {
		return nil, errors.New("分类名称不能为空")
	}

	_, err = l.svcCtx.MenuType.Create(l.ctx, menu.CreateMenuCategoryInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return
}
