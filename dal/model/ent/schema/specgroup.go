package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type SpecGroup struct {
	ent.Schema
}

func (SpecGroup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (SpecGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "spec_group"},
		schema.Comment("规格组"),
	}
}

func (SpecGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(64).
			Comment("规格名（辣度）"),
		field.Int("sort").
			Default(0).
			Comment("排序"),
	}
}

func (SpecGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("spec_items", SpecItem.Type),
	}
}
