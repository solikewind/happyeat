// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateIAMRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建角色（无权限点，需再通过 RBAC 矩阵配置）
func NewCreateIAMRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateIAMRoleLogic {
	return &CreateIAMRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateIAMRoleLogic) CreateIAMRole(req *types.CreateIAMRoleReq) (resp *types.CreateIAMRoleReply, err error) {
	// todo: add your logic here and delete this line

	return
}
