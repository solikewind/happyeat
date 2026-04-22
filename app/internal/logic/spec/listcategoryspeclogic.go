// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	specmodel "github.com/solikewind/happyeat/dal/model/spec"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出分类规格模板
func NewListCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategorySpecLogic {
	return &ListCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCategorySpecLogic) ListCategorySpec(req *types.ListCategorySpecReq) (resp *types.ListCategorySpecReply, err error) {
	offset, limit := normalizePage(req.Current, req.PageSize)
	list, total, err := l.svcCtx.CategorySpec.List(l.ctx, specmodel.ListCategorySpecsFilter{
		CategoryID: req.CategoryId,
		SpecType:   req.SpecType,
		Offset:     offset,
		Limit:      limit,
	})
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &types.ListCategorySpecReply{
			Specs: []types.CategorySpec{},
			Total: 0,
		}, nil
	}

	specs := make([]types.CategorySpec, 0, len(list))
	for _, item := range list {
		specs = append(specs, toCategorySpec(item))
	}

	return &types.ListCategorySpecReply{
		Specs: specs,
		Total: int64(total),
	}, nil
}
