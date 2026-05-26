package order

import (
	"testing"

	"github.com/solikewind/happyeat/dal/model/ent"
)

func ptrU64(v uint64) *uint64 { return &v }
func ptrStr(v string) *string  { return &v }

func mkItem(menuID uint64, name string, qty int, spec string) *ent.OrderItem {
	it := &ent.OrderItem{
		MenuName: name,
		Quantity: qty,
	}
	if menuID > 0 {
		it.MenuID = ptrU64(menuID)
	}
	if spec != "" {
		it.SpecInfo = ptrStr(spec)
	}
	return it
}

func TestItemKey_PrefersMenuID(t *testing.T) {
	a := mkItem(7, "宫保鸡丁", 1, "微辣")
	b := mkItem(7, "宫保鸡", 2, "微辣") // 名字不同但 menu_id 相同
	if itemKey(a) != itemKey(b) {
		t.Fatalf("same menu_id+spec should share key: %q vs %q", itemKey(a), itemKey(b))
	}
}

func TestItemKey_FallbackToName(t *testing.T) {
	a := mkItem(0, "宫保鸡丁", 1, "微辣")
	b := mkItem(0, "宫保鸡丁", 9, "微辣")
	c := mkItem(0, "宫保鸡丁", 1, "中辣")
	if itemKey(a) != itemKey(b) {
		t.Fatalf("same name+spec should share key")
	}
	if itemKey(a) == itemKey(c) {
		t.Fatalf("different spec should differ")
	}
}

func TestDiffOrderItems_AddedRemovedChanged(t *testing.T) {
	old := []*ent.OrderItem{
		mkItem(1, "宫保鸡丁", 2, ""),
		mkItem(2, "鱼香肉丝", 1, ""),
		mkItem(3, "西红柿炒蛋", 1, ""), // 这道菜被删
	}
	now := []*ent.OrderItem{
		mkItem(1, "宫保鸡丁", 1, ""), // 数量从 2 -> 1
		mkItem(2, "鱼香肉丝", 1, ""), // 未变
		mkItem(4, "麻婆豆腐", 1, ""), // 新增
	}
	diff := DiffOrderItems(old, now)

	if len(diff.Removed) != 1 || diff.Removed[0].MenuName != "西红柿炒蛋" {
		t.Fatalf("removed: %+v", diff.Removed)
	}

	added := diff.ByKey[itemKey(now[2])]
	if added.Kind != ItemDiffAdded {
		t.Fatalf("expect added: %+v", added)
	}

	changed := diff.ByKey[itemKey(now[0])]
	if changed.Kind != ItemDiffQtyChanged || changed.OldQty != 2 {
		t.Fatalf("expect qty changed from 2: %+v", changed)
	}

	if _, ok := diff.ByKey[itemKey(now[1])]; ok {
		t.Fatalf("unchanged should not appear in ByKey")
	}
}

func TestDiffOrderItems_EmptySides(t *testing.T) {
	now := []*ent.OrderItem{mkItem(1, "A", 1, "")}
	if d := DiffOrderItems(nil, now); d.ByKey[itemKey(now[0])].Kind != ItemDiffAdded {
		t.Fatalf("nil old -> added")
	}
	old := []*ent.OrderItem{mkItem(1, "A", 1, "")}
	if d := DiffOrderItems(old, nil); len(d.Removed) != 1 {
		t.Fatalf("nil new -> removed")
	}
}
