package order

import (
	"context"
	"errors"
	"strings"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/consts/enum"
	"github.com/solikewind/happyeat/dal/model/ent"
	ordermodel "github.com/solikewind/happyeat/dal/model/order"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新订单（追加菜单项）
func NewUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLogic {
	return &UpdateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrderLogic) UpdateOrder(req *types.UpdateOrderReq) (*types.UpdateOrderReply, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("items 不能为空")
	}
	current, err := l.svcCtx.Order.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}
	if current.Status == enum.OrderStatusCompleted || current.Status == enum.OrderStatusCancelled {
		return nil, errors.New("当前订单状态不允许追加菜单")
	}

	items := make([]ordermodel.ItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		if it.Quantity <= 0 {
			return nil, errors.New("quantity 必须大于 0")
		}

		// 兼容两种入参：
		// 1) 新版：menu_id（推荐）
		// 2) 旧版：menu_name + unit_price（用于历史订单或前端尚未升级时兜底）
		if it.MenuId > 0 {
			menuEnt, err := l.svcCtx.Menu.GetByID(l.ctx, it.MenuId)
			if err != nil {
				if ent.IsNotFound(err) {
					return nil, errors.New("菜单不存在")
				}
				return nil, err
			}
			items = append(items, ordermodel.ItemInput{
				MenuID:    it.MenuId,
				MenuName:  menuEnt.Name,
				Quantity:  it.Quantity,
				UnitPrice: menuEnt.Price,
				SpecInfo:  it.SpecInfo,
			})
			continue
		}

		menuName := strings.TrimSpace(it.MenuName)
		if menuName == "" {
			return nil, errors.New("menu_id 与 menu_name 不能同时为空")
		}
		if it.UnitPrice < 0 {
			return nil, errors.New("unit_price 不能为负")
		}
		items = append(items, ordermodel.ItemInput{
			MenuName:  menuName,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
			SpecInfo:  it.SpecInfo,
		})
	}

	updated, err := l.svcCtx.Order.ReplaceItems(l.ctx, req.Id, items)
	if err != nil {
		return nil, err
	}

	scheduleKitchenPrint(l.svcCtx, updated, "[改单重打]")

	return &types.UpdateOrderReply{Order: EntOrderToType(updated)}, nil
}
