// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除规格组
func NewDeleteSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSpecGroupLogic {
	return &DeleteSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSpecGroupLogic) DeleteSpecGroup(req *types.DeleteSpecGroupReq) (resp *types.DeleteSpecGroupReply, err error) {
	// todo: add your logic here and delete this line

	return
}
