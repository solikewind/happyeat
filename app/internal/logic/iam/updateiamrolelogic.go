// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateIAMRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新角色展示名（role_code 不可改）
func NewUpdateIAMRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateIAMRoleLogic {
	return &UpdateIAMRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateIAMRoleLogic) UpdateIAMRole(req *types.UpdateIAMRoleReq) (resp *types.UpdateIAMRoleReply, err error) {
	// todo: add your logic here and delete this line

	return
}
