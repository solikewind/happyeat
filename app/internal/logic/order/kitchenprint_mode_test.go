package order

import (
	"testing"

	"github.com/solikewind/happyeat/dal/model/ent"
)

func TestResolveKitchenTicketMode(t *testing.T) {
	addDiff := DiffOrderItems(
		[]*ent.OrderItem{mkItem(1, "A", 1, "")},
		[]*ent.OrderItem{mkItem(1, "A", 1, ""), mkItem(2, "B", 1, "")},
	)
	changeDiff := DiffOrderItems(
		[]*ent.OrderItem{mkItem(1, "A", 2, "")},
		[]*ent.OrderItem{mkItem(1, "A", 1, "")},
	)

	if got := resolveKitchenTicketMode("[新单]", addDiff); got != KitchenTicketModeFull {
		t.Fatalf("new order -> full, got %v", got)
	}
	if got := resolveKitchenTicketMode("[手动打印]", addDiff); got != KitchenTicketModeFull {
		t.Fatalf("manual -> full, got %v", got)
	}
	if got := resolveKitchenTicketMode("[改单重打]", addDiff); got != KitchenTicketModeAddOnly {
		t.Fatalf("add-only -> addOnly, got %v", got)
	}
	if got := resolveKitchenTicketMode("[改单重打]", changeDiff); got != KitchenTicketModeChange {
		t.Fatalf("qty change -> change, got %v", got)
	}
}

func TestOrderItemDiff_IsAddOnly(t *testing.T) {
	addOnly := DiffOrderItems(
		[]*ent.OrderItem{mkItem(1, "A", 1, "")},
		[]*ent.OrderItem{mkItem(1, "A", 1, ""), mkItem(2, "B", 1, "")},
	)
	if !addOnly.IsAddOnly() {
		t.Fatal("expected add-only")
	}
	mixed := DiffOrderItems(
		[]*ent.OrderItem{mkItem(1, "A", 2, "")},
		[]*ent.OrderItem{mkItem(1, "A", 1, ""), mkItem(2, "B", 1, "")},
	)
	if mixed.IsAddOnly() {
		t.Fatal("qty change with add is not add-only")
	}
}
