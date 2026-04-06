// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIAMPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个权限点
func NewGetIAMPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIAMPermissionLogic {
	return &GetIAMPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIAMPermissionLogic) GetIAMPermission(req *types.GetIAMPermissionReq) (resp *types.GetIAMPermissionReply, err error) {
	// todo: add your logic here and delete this line

	return
}
