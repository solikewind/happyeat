// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package table

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除餐桌
func NewDeleteTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTableLogic {
	return &DeleteTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTableLogic) DeleteTable(req *types.DeleteTableReq) (resp *types.DeleteTableReply, err error) {
	// todo: add your logic here and delete this line

	return
}
