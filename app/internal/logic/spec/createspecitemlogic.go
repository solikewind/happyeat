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
	name := normalizeText(req.Name)
	if req.SpecGroupId == 0 {
		return nil, errors.New("规格组不能为空")
	}
	if name == "" {
		return nil, errors.New("规格项名称不能为空")
	}

	groupExist, err := l.svcCtx.SpecGroup.Exist(l.ctx, req.SpecGroupId)
	if err != nil {
		return nil, err
	}
	if !groupExist {
		return nil, errors.New("规格组不存在")
	}

	dup, err := l.svcCtx.SpecItem.ExistByName(l.ctx, req.SpecGroupId, name, 0)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, errors.New("规格项已存在")
	}

	if _, err = l.svcCtx.SpecItem.Create(l.ctx, specmodel.CreateSpecItemInput{
		SpecGroupID:  req.SpecGroupId,
		Name:         name,
		DefaultPrice: req.DefaultPrice,
		Sort:         req.Sort,
	}); err != nil {
		return nil, err
	}

	return &types.CreateSpecItemReply{}, nil
}
