package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type MenuSpecOption struct {
	ent.Schema
}

func (MenuSpecOption) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (MenuSpecOption) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "menu_spec_options"},
		schema.Comment("菜单规格选项"),
	}
}

func (MenuSpecOption) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(64).
			Comment("规格选项值"),
		field.Float("price_delta").
			Default(0).
			Comment("加价"),
		field.Int("sort").
			Default(0).
			Comment("排序"),
	}
}

func (MenuSpecOption) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("menu_spec", MenuSpec.Type).Ref("menu_spec_options").Unique().Required(),
	}
}
