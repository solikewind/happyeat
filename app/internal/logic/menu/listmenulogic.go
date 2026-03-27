// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/menu"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

const maxSpecsPerMenuInList = 20

// 列出菜单
func NewListMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuLogic {
	return &ListMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMenuLogic) ListMenu(req *types.ListMenuReq) (resp *types.ListMenuReply, err error) {
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	current := int(req.Current)
	if current <= 0 {
		current = 1
	}
	offset := (current - 1) * pageSize

	list, total, err := l.svcCtx.Menu.List(l.ctx, menu.ListMenusFilter{
		Name:         req.Name,
		CategoryName: req.Category,
		Offset:       offset,
		Limit:        pageSize,
	})
	if err != nil {
		return nil, err
	}

	menus := make([]types.Menu, 0, len(list))
	for _, e := range list {
		menu := entMenuToType(e)
		if len(menu.Specs) > maxSpecsPerMenuInList {
			menu.Specs = menu.Specs[:maxSpecsPerMenuInList]
		}
		menus = append(menus, menu)
	}

	return &types.ListMenuReply{
		Menus: menus,
		Total: total,
	}, nil
}
