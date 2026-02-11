// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tablecategory

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTableCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新餐桌类别
func NewUpdateTableCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTableCategoryLogic {
	return &UpdateTableCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTableCategoryLogic) UpdateTableCategory(req *types.UpdateTableCategoryReq) (resp *types.UpdateTableCategoryReply, err error) {
	// todo: add your logic here and delete this line

	return
}
