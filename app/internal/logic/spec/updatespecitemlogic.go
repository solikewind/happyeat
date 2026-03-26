// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSpecItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新规格项
func NewUpdateSpecItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSpecItemLogic {
	return &UpdateSpecItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSpecItemLogic) UpdateSpecItem(req *types.UpdateSpecItemReq) (resp *types.UpdateSpecItemReply, err error) {
	// todo: add your logic here and delete this line

	return
}
