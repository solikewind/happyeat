// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tablecategory

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/table"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTableCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出餐桌类别
func NewListTableCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTableCategoryLogic {
	return &ListTableCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTableCategoryLogic) ListTableCategory(req *types.ListTableCategoryReq) (resp *types.ListTableCategoryReply, err error) {
	current, pageSize := req.Current, req.PageSize
	if current <= 0 {
		current = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	categories, total, err := l.svcCtx.TableType.List(l.ctx, table.ListTableCategoriesFilter{
		Name:   req.Name,
		Offset: (current - 1) * pageSize,
		Limit:  pageSize,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.TableCategory, 0, len(categories))
	for _, c := range categories {
		list = append(list, entTableCategoryToType(c))
	}
	return &types.ListTableCategoryReply{
		Categories: list,
		Total:      total,
	}, nil
}
