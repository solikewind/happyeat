// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取规格组
func NewGetSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpecGroupLogic {
	return &GetSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSpecGroupLogic) GetSpecGroup(req *types.GetSpecGroupReq) (resp *types.GetSpecGroupReply, err error) {
	// todo: add your logic here and delete this line

	return
}
