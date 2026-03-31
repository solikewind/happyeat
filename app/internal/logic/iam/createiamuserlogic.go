// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateIAMUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建用户（仅主体档案，分配角色用 user-roles 接口）
func NewCreateIAMUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateIAMUserLogic {
	return &CreateIAMUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateIAMUserLogic) CreateIAMUser(req *types.CreateIAMUserReq) (resp *types.CreateIAMUserReply, err error) {
	// todo: add your logic here and delete this line

	return
}
