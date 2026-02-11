// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package table

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTablesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出餐桌
func NewListTablesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTablesLogic {
	return &ListTablesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTablesLogic) ListTables(req *types.ListTablesReq) (resp *types.ListTablesReply, err error) {
	// todo: add your logic here and delete this line

	return
}
