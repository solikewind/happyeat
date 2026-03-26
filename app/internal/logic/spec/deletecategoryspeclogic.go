// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除分类规格模板
func NewDeleteCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategorySpecLogic {
	return &DeleteCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCategorySpecLogic) DeleteCategorySpec(req *types.DeleteCategorySpecReq) (resp *types.DeleteCategorySpecReply, err error) {
	// todo: add your logic here and delete this line

	return
}
