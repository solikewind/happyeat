// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/util/timeutil"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个菜单
func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuLogic) GetMenu(req *types.GetMenuReq) (*types.GetMenuReply, error) {
	entMenu, err := l.svcCtx.Menu.GetByID(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.GetMenuReply{Menu: entMenuToType(entMenu)}, nil
}

func entMenuToType(e *ent.Menu) types.Menu {
	out := types.Menu{
		Id:        uint64(e.ID),
		Name:      e.Name,
		Price:     e.Price,
		CreatedAt: timeutil.TimeToString(e.CreatedAt),
		UpdatedAt: timeutil.TimeToString(e.UpdatedAt),
	}

	if e.Description != nil {
		out.Description = *e.Description
	}
	if e.Image != nil {
		out.Image = *e.Image
	}

	cat, _ := e.Edges.CategoryOrErr()
	if cat != nil {
		out.CategoryId = uint64(cat.ID)
	}

	specs, _ := e.Edges.MenuSpecsOrErr()
	for _, s := range specs {
		spec := types.MenuSpec{
			PriceDelta: s.PriceDelta,
			Sort:       s.Sort,
		}
		if s.SpecItemID != nil {
			spec.SpecItemId = *s.SpecItemID
		}
		if s.CategorySpecID != nil {
			spec.CategorySpecId = *s.CategorySpecID
		}
		out.Specs = append(out.Specs, spec)
	}

	return out
}
