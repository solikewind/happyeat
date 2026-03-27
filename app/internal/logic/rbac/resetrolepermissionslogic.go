// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rbac

import (
	"context"
	"strings"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置角色权限（可单角色）
func NewResetRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetRolePermissionsLogic {
	return &ResetRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetRolePermissionsLogic) ResetRolePermissions(req *types.ResetRolePermissionsReq) (resp *types.ResetRolePermissionsReply, err error) {
	role := strings.TrimSpace(req.Role)
	if err := l.svcCtx.Rbac.Reset(role); err != nil {
		return nil, errInvalid(err.Error())
	}
	if err := svc.SyncRolePoliciesToCasbin(l.svcCtx.Rbac, l.svcCtx.Casbin); err != nil {
		return nil, errInvalid("同步 Casbin 策略失败")
	}
	return &types.ResetRolePermissionsReply{}, nil
}
