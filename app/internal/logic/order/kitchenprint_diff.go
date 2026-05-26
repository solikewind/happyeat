package order

import (
	"fmt"
	"strings"

	"github.com/solikewind/happyeat/dal/model/ent"
)

// ItemDiffKind 菜品变更类型（以「新订单」为视角）。
type ItemDiffKind int

const (
	ItemDiffNone       ItemDiffKind = iota // 未变化（不入 ByKey）
	ItemDiffAdded                          // 新增项
	ItemDiffQtyChanged                     // 同款菜数量发生变化
)

// ItemDiff 单道菜的变更信息。
type ItemDiff struct {
	Kind   ItemDiffKind
	OldQty int // 仅 ItemDiffQtyChanged 时有效
}

// OrderItemDiff 整单 diff 结果。
//
// ByKey: 当前菜单中「有变化」的项；key = itemKey(item)
// Removed: 旧单里存在、新单里不存在的菜（按 sort 顺序）
type OrderItemDiff struct {
	ByKey   map[string]ItemDiff
	Removed []*ent.OrderItem
}

// itemKey 用于识别"同一道菜"的稳定 key。
// 优先 menu_id + spec_info；menu_id=0 时回落到 menu_name + spec_info。
func itemKey(it *ent.OrderItem) string {
	if it == nil {
		return ""
	}
	spec := ""
	if it.SpecInfo != nil {
		spec = strings.TrimSpace(*it.SpecInfo)
	}
	if it.MenuID != nil && *it.MenuID > 0 {
		return fmt.Sprintf("id:%d|%s", *it.MenuID, spec)
	}
	return "name:" + strings.TrimSpace(it.MenuName) + "|" + spec
}

// DiffOrderItems 比对旧菜列表与新菜列表，返回变更结果。
// 任一侧为空也允许：旧空 → 全部 Added；新空 → 全部 Removed。
func DiffOrderItems(oldItems, newItems []*ent.OrderItem) *OrderItemDiff {
	out := &OrderItemDiff{
		ByKey: make(map[string]ItemDiff),
	}

	oldMap := make(map[string]*ent.OrderItem, len(oldItems))
	for _, o := range oldItems {
		if o == nil {
			continue
		}
		oldMap[itemKey(o)] = o
	}

	newMap := make(map[string]*ent.OrderItem, len(newItems))
	for _, n := range newItems {
		if n == nil {
			continue
		}
		k := itemKey(n)
		newMap[k] = n
		if o, ok := oldMap[k]; ok {
			if o.Quantity != n.Quantity {
				out.ByKey[k] = ItemDiff{Kind: ItemDiffQtyChanged, OldQty: o.Quantity}
			}
			// 数量相等则视为未变化，不入 map
		} else {
			out.ByKey[k] = ItemDiff{Kind: ItemDiffAdded}
		}
	}

	for _, o := range oldItems {
		if o == nil {
			continue
		}
		if _, ok := newMap[itemKey(o)]; !ok {
			out.Removed = append(out.Removed, o)
		}
	}

	return out
}
