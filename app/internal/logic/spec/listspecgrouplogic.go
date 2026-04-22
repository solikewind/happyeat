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

type ListSpecGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出规格组
func NewListSpecGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSpecGroupLogic {
	return &ListSpecGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSpecGroupLogic) ListSpecGroup(req *types.ListSpecGroupReq) (resp *types.ListSpecGroupReply, err error) {
	offset, limit := normalizePage(req.Current, req.PageSize)
	list, total, err := l.svcCtx.SpecGroup.List(l.ctx, specmodel.ListSpecGroupsFilter{
		Name:   req.Name,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &types.ListSpecGroupReply{
			Groups: []types.SpecGroup{},
			Total:  0,
		}, nil
	}

	groups := make([]types.SpecGroup, 0, len(list))
	for _, item := range list {
		groups = append(groups, toSpecGroup(item))
	}

	return &types.ListSpecGroupReply{
		Groups: groups,
		Total:  int64(total),
	}, nil
}
