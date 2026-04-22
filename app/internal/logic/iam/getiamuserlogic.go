// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIAMUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个用户及角色列表
func NewGetIAMUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIAMUserLogic {
	return &GetIAMUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIAMUserLogic) GetIAMUser(req *types.GetIAMUserReq) (resp *types.GetIAMUserReply, err error) {
	// todo: add your logic here and delete this line

	return
}
