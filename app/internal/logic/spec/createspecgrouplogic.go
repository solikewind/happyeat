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

type CreateSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建规格组
func NewCreateSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSpecGroupLogic {
	return &CreateSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSpecGroupLogic) CreateSpecGroup(req *types.CreateSpecGroupReq) (resp *types.CreateSpecGroupReply, err error) {
	name := normalizeText(req.Name)
	if name == "" {
		return nil, errors.New("规格组名称不能为空")
	}

	dup, err := l.svcCtx.SpecGroup.ExistByName(l.ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, errors.New("规格组已存在")
	}

	if _, err = l.svcCtx.SpecGroup.Create(l.ctx, specmodel.CreateSpecGroupInput{
		Name: name,
		Sort: req.Sort,
	}); err != nil {
		return nil, err
	}

	return &types.CreateSpecGroupReply{}, nil
}
