// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rbac

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色权限矩阵
func NewListRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolePermissionsLogic {
	return &ListRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRolePermissionsLogic) ListRolePermissions() (resp *types.ListRolePermissionsReply, err error) {
	return &types.ListRolePermissionsReply{
		Roles: l.svcCtx.Rbac.List(),
	}, nil
}
