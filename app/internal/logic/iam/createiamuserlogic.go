// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package iam

import (
	"context"
	"strings"

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
	userCode := strings.TrimSpace(strings.ToLower(req.UserCode))
	displayName := strings.TrimSpace(req.DisplayName)
	id, err := l.svcCtx.Rbac.CreateUser(userCode, displayName)
	if err != nil {
		return nil, errInvalid(err.Error())
	}
	return &types.CreateIAMUserReply{Id: id}, nil
}
