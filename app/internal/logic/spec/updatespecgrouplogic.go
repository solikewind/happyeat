// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新规格组
func NewUpdateSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSpecGroupLogic {
	return &UpdateSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSpecGroupLogic) UpdateSpecGroup(req *types.UpdateSpecGroupReq) (resp *types.UpdateSpecGroupReply, err error) {
	// todo: add your logic here and delete this line

	return
}
