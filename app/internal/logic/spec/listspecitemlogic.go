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

type ListSpecItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列出规格项
func NewListSpecItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSpecItemLogic {
	return &ListSpecItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSpecItemLogic) ListSpecItem(req *types.ListSpecItemReq) (resp *types.ListSpecItemReply, err error) {
	offset, limit := normalizePage(req.Current, req.PageSize)
	list, total, err := l.svcCtx.SpecItem.List(l.ctx, specmodel.ListSpecItemsFilter{
		SpecGroupID: req.SpecGroupId,
		Name:        req.Name,
		Offset:      offset,
		Limit:       limit,
	})
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &types.ListSpecItemReply{
			Items: []types.SpecItem{},
			Total: 0,
		}, nil
	}

	items := make([]types.SpecItem, 0, len(list))
	for _, item := range list {
		items = append(items, toSpecItem(item))
	}

	return &types.ListSpecItemReply{
		Items: items,
		Total: int64(total),
	}, nil
}
