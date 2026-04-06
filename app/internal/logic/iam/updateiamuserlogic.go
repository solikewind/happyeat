// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateIAMUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户展示名
func NewUpdateIAMUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateIAMUserLogic {
	return &UpdateIAMUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateIAMUserLogic) UpdateIAMUser(req *types.UpdateIAMUserReq) (resp *types.UpdateIAMUserReply, err error) {
	// todo: add your logic here and delete this line

	return
}
