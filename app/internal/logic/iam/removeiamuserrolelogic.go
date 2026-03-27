// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveIAMUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 移除用户的某个角色；使用 query：?user_code=&role_code=
func NewRemoveIAMUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveIAMUserRoleLogic {
	return &RemoveIAMUserRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveIAMUserRoleLogic) RemoveIAMUserRole(req *types.RemoveIAMUserRoleReq) (resp *types.RemoveIAMUserRoleReply, err error) {
	// todo: add your logic here and delete this line

	return
}
