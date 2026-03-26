package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type SpecItem struct {
	ent.Schema
}

func (SpecItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (SpecItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "spec_items"},
	}
}

func (SpecItem) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("spec_group_id").
			Comment("所属规格组ID"),
		field.String("name").
			MaxLen(64).
			Comment("规格项名"),
		field.Float("default_price").
			Comment("默认价格"),
	}
}

func (SpecItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("spec_group", SpecGroup.Type).
			Ref("spec_items").
			Field("spec_group_id").
			Unique().
			Required(),
		edge.From("menu_specs", MenuSpec.Type).
			Ref("spec_item"),
	}
}
