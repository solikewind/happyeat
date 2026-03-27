// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListIAMPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页列出权限点（iam_permissions）
func NewListIAMPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListIAMPermissionsLogic {
	return &ListIAMPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListIAMPermissionsLogic) ListIAMPermissions(req *types.ListIAMPermissionsReq) (resp *types.ListIAMPermissionsReply, err error) {
	// todo: add your logic here and delete this line

	return
}
