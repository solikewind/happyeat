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

type UpdateSpecItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新规格项
func NewUpdateSpecItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSpecItemLogic {
	return &UpdateSpecItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSpecItemLogic) UpdateSpecItem(req *types.UpdateSpecItemReq) (resp *types.UpdateSpecItemReply, err error) {
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

	dup, err := l.svcCtx.SpecItem.ExistByName(l.ctx, req.SpecGroupId, name, req.Id)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, errors.New("规格项已存在")
	}

	err = l.svcCtx.SpecItem.Update(l.ctx, req.Id, specmodel.CreateSpecItemInput{
		SpecGroupID:  req.SpecGroupId,
		Name:         name,
		DefaultPrice: req.DefaultPrice,
		Sort:         req.Sort,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("规格项不存在")
		}
		return nil, err
	}

	return &types.UpdateSpecItemReply{}, nil
}
