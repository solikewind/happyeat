package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// MenuCategory holds the schema definition for the MenuCategory entity.
type MenuCategory struct {
	ent.Schema
}

func (MenuCategory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (MenuCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "menu_categories"},
		schema.Comment("菜单分类"),
	}
}

func (MenuCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(64).
			Comment("分类名称"),
		field.String("description").
			Optional().
			Nillable().
			Comment("描述"),
		field.Uint32("sort").
			Default(0).
			Comment("排序，越小越靠前"),
	}
}

func (MenuCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("menus", Menu.Type),
		edge.To("category_specs", CategorySpec.Type),
	}
}
