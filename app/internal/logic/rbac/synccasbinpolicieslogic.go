// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rbac

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncCasbinPoliciesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 将 IAM 同步到 Casbin（刷新 casbin_rule，供管理端按钮触发）
func NewSyncCasbinPoliciesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncCasbinPoliciesLogic {	return &SyncCasbinPoliciesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncCasbinPoliciesLogic) SyncCasbinPolicies(req *types.SyncCasbinPoliciesReq) (resp *types.SyncCasbinPoliciesReply, err error) {
	_ = req
	if syncErr := svc.SyncRolePoliciesToCasbin(l.svcCtx.Rbac, l.svcCtx.Casbin); syncErr != nil {
		l.Errorf("sync casbin: %v", syncErr)
		return nil, errInvalid("同步 Casbin 策略失败")
	}
	return &types.SyncCasbinPoliciesReply{}, nil
}
