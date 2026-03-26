// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建分类规格模板
func NewCreateCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategorySpecLogic {
	return &CreateCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategorySpecLogic) CreateCategorySpec(req *types.CreateCategorySpecReq) (resp *types.CreateCategorySpecReply, err error) {
	// todo: add your logic here and delete this line

	return
}
