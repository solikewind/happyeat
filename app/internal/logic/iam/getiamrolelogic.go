// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIAMRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个角色
func NewGetIAMRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIAMRoleLogic {
	return &GetIAMRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIAMRoleLogic) GetIAMRole(req *types.GetIAMRoleReq) (resp *types.GetIAMRoleReply, err error) {
	// todo: add your logic here and delete this line

	return
}
