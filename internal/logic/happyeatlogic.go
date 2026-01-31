// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"happyeat/internal/svc"
	"happyeat/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HappyeatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHappyeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HappyeatLogic {
	return &HappyeatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HappyeatLogic) Happyeat(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
