// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tablecategory

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTableCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出餐桌类别
func NewListTableCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTableCategoriesLogic {
	return &ListTableCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTableCategoriesLogic) ListTableCategories(req *types.ListTableCategoriesReq) (resp *types.ListTableCategoriesReply, err error) {
	// todo: add your logic here and delete this line

	return
}
