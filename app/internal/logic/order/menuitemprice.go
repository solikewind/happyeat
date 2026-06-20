package order

import (
	"strings"

	"github.com/solikewind/happyeat/dal/model/ent"
)

type menuSpecKey struct {
	typ string
	val string
}

// resolveMenuUnitPrice 根据菜单基础价、规格描述与客户端单价（可选）计算行单价。
func resolveMenuUnitPrice(menuEnt *ent.Menu, specInfo string, clientUnitPrice int64) int64 {
	base := menuEnt.Price
	if clientUnitPrice > 0 {
		return clientUnitPrice
	}
	return base + specPriceDelta(menuEnt, specInfo)
}

func specPriceDelta(menuEnt *ent.Menu, specInfo string) int64 {
	specInfo = strings.TrimSpace(specInfo)
	if specInfo == "" {
		return 0
	}
	specs := menuSpecList(menuEnt)
	if len(specs) == 0 {
		return 0
	}

	lookup := make(map[menuSpecKey]int64, len(specs))
	for _, s := range specs {
		typ, val := menuSpecTypeValue(s)
		if typ == "" || val == "" {
			continue
		}
		lookup[menuSpecKey{typ: typ, val: val}] += menuSpecPriceDelta(s)
	}

	var delta int64
	for _, part := range strings.Fields(specInfo) {
		typ, val, ok := strings.Cut(part, ":")
		if !ok || typ == "" || val == "" {
			continue
		}
		delta += lookup[menuSpecKey{typ: typ, val: val}]
	}
	return delta
}

func menuSpecPriceDelta(s *ent.MenuSpec) int64 {
	if categorySpec, err := s.Edges.CategorySpecOrErr(); err == nil && categorySpec != nil {
		return categorySpecCurrentPriceDelta(categorySpec)
	}
	if s.Edges.CategorySpec != nil {
		return categorySpecCurrentPriceDelta(s.Edges.CategorySpec)
	}
	if specItem, err := s.Edges.SpecItemOrErr(); err == nil && specItem != nil {
		return specItem.DefaultPrice
	}
	if s.Edges.SpecItem != nil {
		return s.Edges.SpecItem.DefaultPrice
	}
	return s.PriceDelta
}

func categorySpecCurrentPriceDelta(categorySpec *ent.CategorySpec) int64 {
	if specItem, err := categorySpec.Edges.SpecItemOrErr(); err == nil && specItem != nil {
		return specItem.DefaultPrice
	}
	if categorySpec.Edges.SpecItem != nil {
		return categorySpec.Edges.SpecItem.DefaultPrice
	}
	return categorySpec.PriceDelta
}

func menuSpecList(menuEnt *ent.Menu) []*ent.MenuSpec {
	if specs, err := menuEnt.Edges.MenuSpecsOrErr(); err == nil {
		return specs
	}
	return menuEnt.Edges.MenuSpecs
}

func menuSpecTypeValue(s *ent.MenuSpec) (typ, val string) {
	var categorySpec *ent.CategorySpec
	if cs, err := s.Edges.CategorySpecOrErr(); err == nil {
		categorySpec = cs
	} else if s.Edges.CategorySpec != nil {
		categorySpec = s.Edges.CategorySpec
	}
	if categorySpec != nil {
		typ = categorySpec.SpecType
		val = categorySpec.SpecValue
	}

	var specItem *ent.SpecItem
	if si, err := s.Edges.SpecItemOrErr(); err == nil {
		specItem = si
	} else if s.Edges.SpecItem != nil {
		specItem = s.Edges.SpecItem
	}
	if specItem != nil {
		if val == "" {
			val = specItem.Name
		}
		if typ == "" {
			var group *ent.SpecGroup
			if g, err := specItem.Edges.SpecGroupOrErr(); err == nil {
				group = g
			} else if specItem.Edges.SpecGroup != nil {
				group = specItem.Edges.SpecGroup
			}
			if group != nil {
				typ = group.Name
			}
		}
	}
	return typ, val
}
