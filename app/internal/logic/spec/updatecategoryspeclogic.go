// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新分类规格模板
func NewUpdateCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCategorySpecLogic {
	return &UpdateCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategorySpecLogic) UpdateCategorySpec(req *types.UpdateCategorySpecReq) (resp *types.UpdateCategorySpecReply, err error) {
	// todo: add your logic here and delete this line

	return
}
