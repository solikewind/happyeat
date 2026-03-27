// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListIAMRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页列出角色（iam_roles）
func NewListIAMRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListIAMRolesLogic {
	return &ListIAMRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListIAMRolesLogic) ListIAMRoles(req *types.ListIAMRolesReq) (resp *types.ListIAMRolesReply, err error) {
	// todo: add your logic here and delete this line

	return
}
