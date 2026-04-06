package menucategory

import (
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/util/timeutil"
	"github.com/solikewind/happyeat/dal/model/ent"
)

func entMenuCategoryToType(e *ent.MenuCategory) types.MenuCategory {
	out := types.MenuCategory{
		Id:        uint64(e.ID),
		Name:      e.Name,
		CreatedAt: timeutil.TimeToString(e.CreatedAt),
		UpdatedAt: timeutil.TimeToString(e.UpdatedAt),
	}
	if e.Description != nil {
		out.Description = *e.Description
	}
	return out
}
