// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package table

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建餐桌
func NewCreateTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTableLogic {
	return &CreateTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTableLogic) CreateTable(req *types.CreateTableReq) (resp *types.CreateTableReply, err error) {

	return
}
