// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menucategory

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
