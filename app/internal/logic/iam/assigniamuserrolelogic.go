// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignIAMUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 为用户绑定角色（幂等：已绑定则成功）；path 固定便于 Casbin obj 精确匹配
func NewAssignIAMUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignIAMUserRoleLogic {
	return &AssignIAMUserRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignIAMUserRoleLogic) AssignIAMUserRole(req *types.AssignIAMUserRoleReq) (resp *types.AssignIAMUserRoleReply, err error) {
	// todo: add your logic here and delete this line

	return
}
