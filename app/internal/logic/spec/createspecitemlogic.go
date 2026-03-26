// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSpecItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建规格项
func NewCreateSpecItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSpecItemLogic {
	return &CreateSpecItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSpecItemLogic) CreateSpecItem(req *types.CreateSpecItemReq) (resp *types.CreateSpecItemReply, err error) {
	// todo: add your logic here and delete this line

	return
}
