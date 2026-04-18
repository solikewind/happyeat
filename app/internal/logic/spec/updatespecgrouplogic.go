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

type UpdateSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新规格组
func NewUpdateSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSpecGroupLogic {
	return &UpdateSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSpecGroupLogic) UpdateSpecGroup(req *types.UpdateSpecGroupReq) (resp *types.UpdateSpecGroupReply, err error) {
	name := normalizeText(req.Name)
	if name == "" {
		return nil, errors.New("规格组名称不能为空")
	}

	dup, err := l.svcCtx.SpecGroup.ExistByName(l.ctx, name, req.Id)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, errors.New("规格组已存在")
	}

	err = l.svcCtx.SpecGroup.Update(l.ctx, req.Id, specmodel.CreateSpecGroupInput{
		Name: name,
		Sort: req.Sort,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("规格组不存在")
		}
		return nil, err
	}

	return &types.UpdateSpecGroupReply{}, nil
}
