// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分类规格模板
func NewGetCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategorySpecLogic {
	return &GetCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategorySpecLogic) GetCategorySpec(req *types.GetCategorySpecReq) (resp *types.GetCategorySpecReply, err error) {
	item, err := l.svcCtx.CategorySpec.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("分类规格不存在")
		}
		return nil, err
	}

	return &types.GetCategorySpecReply{
		Spec: toCategorySpec(item),
	}, nil
}
