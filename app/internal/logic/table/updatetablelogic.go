// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package table

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新餐桌
func NewUpdateTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTableLogic {
	return &UpdateTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTableLogic) UpdateTable(req *types.UpdateTableReq) (resp *types.UpdateTableReply, err error) {
	// todo: add your logic here and delete this line

	return
}
