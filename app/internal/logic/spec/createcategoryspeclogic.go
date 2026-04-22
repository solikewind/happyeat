// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	specmodel "github.com/solikewind/happyeat/dal/model/spec"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建分类规格模板
func NewCreateCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategorySpecLogic {
	return &CreateCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategorySpecLogic) CreateCategorySpec(req *types.CreateCategorySpecReq) (resp *types.CreateCategorySpecReply, err error) {
	if req.CategoryId == 0 {
		return nil, errors.New("分类不能为空")
	}

	exist, err := l.svcCtx.MenuType.Exist(l.ctx, req.CategoryId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("分类不存在")
	}

	specItemID, specType, specValue, err := resolveCategorySpecSource(l.ctx, l.svcCtx, req.SpecItemId, req.SpecType, req.SpecValue, req.PriceDelta)
	if err != nil {
		return nil, err
	}

	dup, err := l.svcCtx.CategorySpec.ExistByValue(l.ctx, req.CategoryId, specType, specValue, 0)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, errors.New("分类规格已存在")
	}

	if _, err = l.svcCtx.CategorySpec.Create(l.ctx, specmodel.CreateCategorySpecInput{
		CategoryID: req.CategoryId,
		SpecItemID: specItemID,
		SpecType:   specType,
		SpecValue:  specValue,
		PriceDelta: req.PriceDelta,
		Sort:       req.Sort,
	}); err != nil {
		return nil, err
	}

	return &types.CreateCategorySpecReply{}, nil
}
