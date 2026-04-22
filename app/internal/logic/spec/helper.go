package spec

import (
	"strings"

	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/util/timeutil"
	"github.com/solikewind/happyeat/dal/model/ent"
)

const defaultPageSize = 10

func normalizePage(current, pageSize int64) (offset, limit int) {
	limit = int(pageSize)
	if limit <= 0 {
		limit = defaultPageSize
	}

	page := int(current)
	if page <= 0 {
		page = 1
	}

	return (page - 1) * limit, limit
}

func normalizeText(v string) string {
	return strings.TrimSpace(v)
}

func toCategorySpec(e *ent.CategorySpec) types.CategorySpec {
	out := types.CategorySpec{
		Id:         e.ID,
		CategoryId: e.MenuCategoryID,
		SpecType:   e.SpecType,
		SpecValue:  e.SpecValue,
		PriceDelta: e.PriceDelta,
		Sort:       e.Sort,
		CreatedAt:  timeutil.TimeToString(e.CreatedAt),
		UpdatedAt:  timeutil.TimeToString(e.UpdatedAt),
	}
	if e.SpecItemID != nil {
		out.SpecItemId = *e.SpecItemID
	}
	if item, err := e.Edges.SpecItemOrErr(); err == nil && item != nil {
		out.SpecItemId = item.ID
		out.SpecValue = item.Name
		if group, groupErr := item.Edges.SpecGroupOrErr(); groupErr == nil && group != nil {
			out.SpecType = group.Name
		}
	}
	return out
}

func toSpecGroup(e *ent.SpecGroup) types.SpecGroup {
	return types.SpecGroup{
		Id:        e.ID,
		Name:      e.Name,
		Sort:      e.Sort,
		CreatedAt: timeutil.TimeToString(e.CreatedAt),
		UpdatedAt: timeutil.TimeToString(e.UpdatedAt),
	}
}

func toSpecItem(e *ent.SpecItem) types.SpecItem {
	return types.SpecItem{
		Id:           e.ID,
		SpecGroupId:  e.SpecGroupID,
		Name:         e.Name,
		DefaultPrice: e.DefaultPrice,
		CreatedAt:    timeutil.TimeToString(e.CreatedAt),
		UpdatedAt:    timeutil.TimeToString(e.UpdatedAt),
	}
}
