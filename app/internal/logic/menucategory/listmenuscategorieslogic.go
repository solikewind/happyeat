// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package MenuCategory

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/util/paging"
	menudata "github.com/solikewind/happyeat/dal/model/menu"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenusCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出菜单种类
func NewListMenusCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenusCategoriesLogic {
	return &ListMenusCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMenusCategoriesLogic) ListMenusCategories(req *types.ListMenusCategoriesReq) (*types.ListMenusCategoriesReply, error) {
	pageParam := paging.NewPageParam(req.Current, req.PageSize)
	offset := pageParam.Offset()
	pageSize := pageParam.PageSize

	list, total, err := l.svcCtx.MenuType.List(l.ctx, menudata.ListMenuCategoriesFilter{
		Name:   req.Name,
		Offset: offset,
		Limit:  pageSize,
	})
	if err != nil {
		return nil, err
	}

	categories := make([]types.MenuCategory, 0, len(list))
	for _, e := range list {
		categories = append(categories, entCategoryToType(e))
	}

	return &types.ListMenusCategoriesReply{
		Categories: categories,
		Total:      total,
	}, nil
}
