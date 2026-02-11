// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/menu/ent"

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
	entMenu, err := l.svcCtx.Menu.GetByID(l.ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &types.GetMenuReply{Menu: entMenuToType(entMenu)}, nil
}

func entMenuToType(e *ent.Menu) types.Menu {
	out := types.Menu{
		Id:       uint64(e.ID),
		Name:     e.Name,
		Price:    e.Price,
		CreateAt: e.CreatedAt.Unix(),
		UpdateAt: e.UpdatedAt.Unix(),
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

	specs, _ := e.Edges.SpecsOrErr()
	for _, s := range specs {
		out.Specs = append(out.Specs, types.MenuSpec{
			SpecType:   s.SpecType,
			SpecValue:  s.SpecValue,
			PriceDelta: s.PriceDelta,
		})
	}

	return out
}
