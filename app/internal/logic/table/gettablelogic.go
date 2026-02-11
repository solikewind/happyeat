// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package table

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个餐桌
func NewGetTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTableLogic {
	return &GetTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTableLogic) GetTable(req *types.GetTableReq) (resp *types.GetTableReply, err error) {
	// todo: add your logic here and delete this line

	return
}
