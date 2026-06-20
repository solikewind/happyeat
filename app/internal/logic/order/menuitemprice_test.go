package order

import (
	"testing"

	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/ent/categoryspec"
	"github.com/solikewind/happyeat/dal/model/ent/menuspec"
)

func TestResolveMenuUnitPrice_WithClientPrice(t *testing.T) {
	menu := &ent.Menu{Price: 28}
	got := resolveMenuUnitPrice(menu, "辣度:微辣", 29)
	if got != 29 {
		t.Fatalf("client unit price = %d, want 29", got)
	}
}

func TestResolveMenuUnitPrice_WithSpecDelta(t *testing.T) {
	menu := &ent.Menu{
		Price: 28,
		Edges: ent.MenuEdges{
			MenuSpecs: []*ent.MenuSpec{
				{
					PriceDelta: 1,
					Edges: ent.MenuSpecEdges{
						CategorySpec: &ent.CategorySpec{
							SpecType:   "辣度",
							SpecValue:  "微辣",
							PriceDelta: 1,
						},
					},
				},
			},
		},
	}
	got := resolveMenuUnitPrice(menu, "辣度:微辣", 0)
	if got != 29 {
		t.Fatalf("resolved = %d, want 29", got)
	}
}

func TestSpecPriceDelta_MultipleSpecs(t *testing.T) {
	menu := &ent.Menu{
		Edges: ent.MenuEdges{
			MenuSpecs: []*ent.MenuSpec{
				{
					PriceDelta: 1,
					Edges: ent.MenuSpecEdges{
						CategorySpec: &ent.CategorySpec{SpecType: "辣度", SpecValue: "重辣", PriceDelta: 1},
					},
				},
				{
					PriceDelta: 2,
					Edges: ent.MenuSpecEdges{
						CategorySpec: &ent.CategorySpec{SpecType: "份量", SpecValue: "大份", PriceDelta: 2},
					},
				},
			},
		},
	}
	delta := specPriceDelta(menu, "辣度:重辣 份量:大份")
	if delta != 3 {
		t.Fatalf("delta = %d, want 3", delta)
	}
}

func TestSpecPriceDelta_UsesCategorySpecCurrentPrice(t *testing.T) {
	menu := &ent.Menu{
		Edges: ent.MenuEdges{
			MenuSpecs: []*ent.MenuSpec{
				{
					PriceDelta: 1,
					Edges: ent.MenuSpecEdges{
						CategorySpec: &ent.CategorySpec{SpecType: "份量", SpecValue: "大份", PriceDelta: 2},
					},
				},
			},
		},
	}
	delta := specPriceDelta(menu, "份量:大份")
	if delta != 2 {
		t.Fatalf("delta = %d, want 2", delta)
	}
}

func TestSpecPriceDelta_UsesSpecItemDefaultPriceThroughCategorySpec(t *testing.T) {
	menu := &ent.Menu{
		Edges: ent.MenuEdges{
			MenuSpecs: []*ent.MenuSpec{
				{
					PriceDelta: 1,
					Edges: ent.MenuSpecEdges{
						CategorySpec: &ent.CategorySpec{
							SpecType:   "份量",
							SpecValue:  "大份",
							PriceDelta: 1,
							Edges: ent.CategorySpecEdges{
								SpecItem: &ent.SpecItem{DefaultPrice: 2},
							},
						},
					},
				},
			},
		},
	}
	delta := specPriceDelta(menu, "份量:大份")
	if delta != 2 {
		t.Fatalf("delta = %d, want 2", delta)
	}
}

func TestMenuSpecTypeValue_FromSpecItem(t *testing.T) {
	s := &ent.MenuSpec{
		PriceDelta: 1,
		Edges: ent.MenuSpecEdges{
			SpecItem: &ent.SpecItem{
				Name: "大杯",
				Edges: ent.SpecItemEdges{
					SpecGroup: &ent.SpecGroup{Name: "容量"},
				},
			},
		},
	}
	typ, val := menuSpecTypeValue(s)
	if typ != "容量" || val != "大杯" {
		t.Fatalf("got %q:%q, want 容量:大杯", typ, val)
	}
	_ = menuspec.Table
	_ = categoryspec.Table
}
