package settlement

import (
	"context"

	orderlogic "github.com/solikewind/happyeat/app/internal/logic/order"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/util/timeutil"
	"github.com/solikewind/happyeat/dal/model/ent"
)

func EntToType(ctx context.Context, svcCtx *svc.ServiceContext, e *ent.Settlement, withOrders bool) types.Settlement {
	out := types.Settlement{
		Id:           uint64(e.ID),
		CustomerName: e.CustomerName,
		Status:       string(e.Status),
		TotalAmount:  e.TotalAmount,
		ActualAmount: e.ActualAmount,
		CreatedAt:    timeutil.TimeToString(e.CreatedAt),
		UpdatedAt:    timeutil.TimeToString(e.UpdatedAt),
	}
	if e.Remark != nil {
		out.Remark = *e.Remark
	}
	if e.SettledAt != nil {
		out.SettledAt = timeutil.TimeToString(*e.SettledAt)
	}
	if withOrders {
		orders, _ := e.Edges.OrdersOrErr()
		out.OrderCount = len(orders)
		for _, ord := range orders {
			out.Orders = append(out.Orders, orderlogic.EntOrderToTypeForDisplay(ctx, svcCtx, ord))
		}
	}
	return out
}

func EntListItemToType(e *ent.Settlement, orderCount int) types.Settlement {
	out := types.Settlement{
		Id:           uint64(e.ID),
		CustomerName: e.CustomerName,
		Status:       string(e.Status),
		TotalAmount:  e.TotalAmount,
		ActualAmount: e.ActualAmount,
		OrderCount:   orderCount,
		CreatedAt:    timeutil.TimeToString(e.CreatedAt),
		UpdatedAt:    timeutil.TimeToString(e.UpdatedAt),
	}
	if e.Remark != nil {
		out.Remark = *e.Remark
	}
	if e.SettledAt != nil {
		out.SettledAt = timeutil.TimeToString(*e.SettledAt)
	}
	return out
}
