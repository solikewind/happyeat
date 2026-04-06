// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menucategory

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/menu"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenuCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出菜单种类
func NewListMenuCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuCategoryLogic {
	return &ListMenuCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMenuCategoryLogic) ListMenuCategory(req *types.ListMenuCategoryReq) (resp *types.ListMenuCategoryReply, err error) {
	current, pageSize := req.Current, req.PageSize
	if current <= 0 {
		current = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := int((current - 1) * pageSize)
	limit := int(pageSize)

	categories, total, err := l.svcCtx.MenuType.List(l.ctx, menu.ListMenuCategoriesFilter{
		Name:   req.Name,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		l.Errorf("ListMenuCategory MenuType.List err: %v", err)
		return nil, err
	}
	list := make([]types.MenuCategory, 0, len(categories))
	for _, c := range categories {
		list = append(list, entMenuCategoryToType(c))
	}
	return &types.ListMenuCategoryReply{
		Categories: list,
		Total:      total,
	}, nil
}
