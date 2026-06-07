package order

import (
	"context"
	"sort"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/menu"
)

type orderItemSortKey struct {
	kindRank     int
	categorySort uint32
	menuSort     uint32
	itemSort     uint32
	id           uint64
}

// ApplyOrderItemsDisplaySort 按菜单分类重排订单明细：菜品类靠前，酒水饮料类靠后。
func ApplyOrderItemsDisplaySort(ctx context.Context, svcCtx *svc.ServiceContext, e *ent.Order) {
	if svcCtx == nil || e == nil {
		return
	}
	items, err := e.Edges.ItemsOrErr()
	if err != nil || len(items) <= 1 {
		return
	}

	ranked, err := svcCtx.Menu.ItemDisplayRanksByMenuIDs(ctx, collectMenuIDsFromItems(items))
	if err != nil {
		return
	}
	e.Edges.Items = sortOrderItemsForDisplay(items, ranked)
}

// EntOrderToTypeForDisplay 先按分类排序再转为 API 类型。
func EntOrderToTypeForDisplay(ctx context.Context, svcCtx *svc.ServiceContext, e *ent.Order) types.Order {
	ApplyOrderItemsDisplaySort(ctx, svcCtx, e)
	return EntOrderToType(e)
}

func collectMenuIDsFromItems(items []*ent.OrderItem) []uint64 {
	ids := make([]uint64, 0, len(items))
	for _, it := range items {
		if it == nil || it.MenuID == nil || *it.MenuID == 0 {
			continue
		}
		ids = append(ids, *it.MenuID)
	}
	return ids
}

func sortOrderItemsForDisplay(items []*ent.OrderItem, ranked map[uint64]menu.ItemDisplayRank) []*ent.OrderItem {
	out := make([]*ent.OrderItem, len(items))
	copy(out, items)
	sort.SliceStable(out, func(i, j int) bool {
		ki := orderItemSortKeyFor(out[i], ranked)
		kj := orderItemSortKeyFor(out[j], ranked)
		if ki.kindRank != kj.kindRank {
			return ki.kindRank < kj.kindRank
		}
		if ki.categorySort != kj.categorySort {
			return ki.categorySort < kj.categorySort
		}
		if ki.menuSort != kj.menuSort {
			return ki.menuSort < kj.menuSort
		}
		if ki.itemSort != kj.itemSort {
			return ki.itemSort < kj.itemSort
		}
		return ki.id < kj.id
	})
	return out
}

func orderItemSortKeyFor(it *ent.OrderItem, ranked map[uint64]menu.ItemDisplayRank) orderItemSortKey {
	key := orderItemSortKey{itemSort: it.Sort, id: uint64(it.ID)}
	if it.MenuID != nil && *it.MenuID > 0 {
		if r, ok := ranked[*it.MenuID]; ok {
			key.kindRank = r.KindRank
			key.categorySort = r.CategorySort
			key.menuSort = r.MenuSort
		}
	}
	return key
}
