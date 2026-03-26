// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建规格组
func NewCreateSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSpecGroupLogic {
	return &CreateSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSpecGroupLogic) CreateSpecGroup(req *types.CreateSpecGroupReq) (resp *types.CreateSpecGroupReply, err error) {
	// todo: add your logic here and delete this line

	return
}
