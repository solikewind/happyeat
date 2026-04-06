// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteIAMRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除角色（软删；系统预置角色不可删；删除后全量同步 Casbin）
func NewDeleteIAMRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteIAMRoleLogic {
	return &DeleteIAMRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteIAMRoleLogic) DeleteIAMRole(req *types.DeleteIAMRoleReq) (resp *types.DeleteIAMRoleReply, err error) {
	// todo: add your logic here and delete this line

	return
}
