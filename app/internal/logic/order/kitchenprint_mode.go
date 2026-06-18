package order

import "strings"

// KitchenTicketMode 厨房小票排版模式。
type KitchenTicketMode int

const (
	// KitchenTicketModeFull 完整单：新单、手动补打。
	KitchenTicketModeFull KitchenTicketMode = iota
	// KitchenTicketModeAddOnly 纯加菜增量单：仅新增项，无整单合计。
	KitchenTicketModeAddOnly
	// KitchenTicketModeChange 变更单：含删菜或改量。
	KitchenTicketModeChange
)

func resolveKitchenTicketMode(banner string, diff *OrderItemDiff) KitchenTicketMode {
	trimmed := strings.Trim(strings.TrimSpace(banner), "[]【】")
	switch trimmed {
	case "新单", "手动打印":
		return KitchenTicketModeFull
	case "改单重打":
		if diff == nil || !diff.HasChanges() {
			return KitchenTicketModeFull
		}
		if diff.IsAddOnly() {
			return KitchenTicketModeAddOnly
		}
		return KitchenTicketModeChange
	default:
		return KitchenTicketModeFull
	}
}

func (d *OrderItemDiff) HasChanges() bool {
	if d == nil {
		return false
	}
	return len(d.ByKey) > 0 || len(d.Removed) > 0
}
