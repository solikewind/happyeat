// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出分类规格模板
func NewListCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategorySpecLogic {
	return &ListCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCategorySpecLogic) ListCategorySpec(req *types.ListCategorySpecReq) (resp *types.ListCategorySpecReply, err error) {
	// todo: add your logic here and delete this line

	return
}
