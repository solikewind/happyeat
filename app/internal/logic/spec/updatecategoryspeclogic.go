// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"
	specmodel "github.com/solikewind/happyeat/dal/model/spec"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategorySpecLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新分类规格模板
func NewUpdateCategorySpecLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCategorySpecLogic {
	return &UpdateCategorySpecLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategorySpecLogic) UpdateCategorySpec(req *types.UpdateCategorySpecReq) (resp *types.UpdateCategorySpecReply, err error) {
	if req.Id == 0 {
		return nil, errors.New("规格ID不能为空")
	}
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

	dup, err := l.svcCtx.CategorySpec.ExistByValue(l.ctx, req.CategoryId, specType, specValue, req.Id)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, errors.New("分类规格已存在")
	}

	err = l.svcCtx.CategorySpec.Update(l.ctx, req.Id, specmodel.CreateCategorySpecInput{
		CategoryID: req.CategoryId,
		SpecItemID: specItemID,
		SpecType:   specType,
		SpecValue:  specValue,
		PriceDelta: req.PriceDelta,
		Sort:       req.Sort,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("分类规格不存在")
		}
		return nil, err
	}

	return &types.UpdateCategorySpecReply{}, nil
}
