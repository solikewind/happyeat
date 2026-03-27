// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出规格组
func NewListSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSpecGroupLogic {
	return &ListSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSpecGroupLogic) ListSpecGroup(req *types.ListSpecGroupReq) (resp *types.ListSpecGroupReply, err error) {
	// todo: add your logic here and delete this line

	return
}
