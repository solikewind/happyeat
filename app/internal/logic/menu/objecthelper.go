package menu

import (
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/util/timeutil"
	"github.com/solikewind/happyeat/dal/model/ent"
)

func entObjectToType(e *ent.Object) types.Object {
	out := types.Object{
		Id:        e.ID,
		Name:      e.Name,
		Key:       e.Key,
		Url:       e.URL,
		Size:      e.Size,
		Hash:      e.Hash,
		CreatedAt: timeutil.TimeToString(e.CreatedAt),
		UpdatedAt: timeutil.TimeToString(e.UpdatedAt),
	}
	if e.ContentType != nil {
		out.ContentType = *e.ContentType
	}
	if e.Suffix != nil {
		out.Suffix = *e.Suffix
	}
	return out
}
