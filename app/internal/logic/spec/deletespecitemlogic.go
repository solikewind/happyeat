// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSpecItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除规格项
func NewDeleteSpecItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSpecItemLogic {
	return &DeleteSpecItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSpecItemLogic) DeleteSpecItem(req *types.DeleteSpecItemReq) (resp *types.DeleteSpecItemReply, err error) {
	// todo: add your logic here and delete this line

	return
}
