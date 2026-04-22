// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteIAMUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除用户（软删，并清除 Casbin 中该用户全部分组）
func NewDeleteIAMUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteIAMUserLogic {
	return &DeleteIAMUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteIAMUserLogic) DeleteIAMUser(req *types.DeleteIAMUserReq) (resp *types.DeleteIAMUserReply, err error) {
	// todo: add your logic here and delete this line

	return
}
