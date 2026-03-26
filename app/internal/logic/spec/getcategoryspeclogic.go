// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分类规格模板
func NewGetCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategorySpecLogic {
	return &GetCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategorySpecLogic) GetCategorySpec(req *types.GetCategorySpecReq) (resp *types.GetCategorySpecReply, err error) {
	// todo: add your logic here and delete this line

	return
}
