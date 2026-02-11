// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tablecategory

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTableCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除餐桌类别
func NewDeleteTableCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTableCategoryLogic {
	return &DeleteTableCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTableCategoryLogic) DeleteTableCategory(req *types.DeleteTableCategoryReq) (resp *types.DeleteTableCategoryReply, err error) {
	// todo: add your logic here and delete this line

	return
}
