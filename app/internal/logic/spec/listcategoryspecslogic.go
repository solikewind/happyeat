// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategorySpecsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出分类规格模板
func NewListCategorySpecsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategorySpecsLogic {
	return &ListCategorySpecsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCategorySpecsLogic) ListCategorySpecs(req *types.ListCategorySpecsReq) (resp *types.ListCategorySpecsReply, err error) {
	// todo: add your logic here and delete this line

	return
}
