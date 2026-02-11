// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tablecategory

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTableCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建餐桌类别
func NewCreateTableCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTableCategoryLogic {
	return &CreateTableCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTableCategoryLogic) CreateTableCategory(req *types.CreateTableCategoryReq) (resp *types.CreateTableCategoryReply, err error) {
	// todo: add your logic here and delete this line

	return
}
