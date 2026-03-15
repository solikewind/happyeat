package schema

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	gen "github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/ent/hook"
	"github.com/solikewind/happyeat/dal/model/ent/intercept"
)

var CST = time.FixedZone("CST", 8*3600)

func CSTNow() time.Time {
	return time.Now().In(CST)
}

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(CSTNow).
			Immutable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}).
			Comment("创建时间"),
		field.Time("updated_at").
			Default(CSTNow).
			UpdateDefault(CSTNow).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}).
			Comment("更新时间"),
	}
}

// SoftDeleteMixin implements the soft delete pattern for schemas.
type SoftDeleteMixin struct {
	mixin.Schema
}

// Fields of the SoftDeleteMixin.
func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("delete_ts").
			Default(0).
			Comment("删除时间戳"),
	}
}

type softDeleteKey struct{}

// SkipSoftDelete returns a new context that skips the soft-delete interceptor/mutators.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// Interceptors of the SoftDeleteMixin.
func (d SoftDeleteMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
			// Skip soft-delete, means include soft-deleted entities.
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}
			d.P(q)
			return nil
		}),
	}
}

// Hooks of the SoftDelete.
func (d SoftDeleteMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
					// Skip soft-delete, means delete the entity permanently.
					if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
						return next.Mutate(ctx, m)
					}
					mx, ok := m.(interface {
						SetOp(ent.Op)
						Client() *gen.Client
						SetDeleteTs(int64)
						WhereP(...func(*sql.Selector))
					})
					if !ok {
						return nil, fmt.Errorf("unexpected mutation type %T", m)
					}
					d.P(mx)
					mx.SetOp(ent.OpUpdate)
					mx.SetDeleteTs(CSTNow().Unix())
					return mx.Client().Mutate(ctx, m)
				})
			},
			ent.OpDeleteOne|ent.OpDelete,
		),
	}
}

// P adds a storage-level predicate to the queries and mutations.
func (d SoftDeleteMixin) P(w interface{ WhereP(...func(*sql.Selector)) }) {
	w.WhereP(
		sql.FieldEQ(d.Fields()[0].Descriptor().Name, int64(0)),
	)
}
