package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// MenuSpec holds the schema definition for the MenuSpec entity.
type MenuSpec struct {
	ent.Schema
}

func (MenuSpec) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (MenuSpec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "menu_specs"},
		schema.Comment("菜单规格（如口味、大小及加价）"),
	}
}

func (MenuSpec) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("menu_id").
			Comment("菜单ID"),
		field.Uint64("spec_item_id").
			Optional().
			Nillable().
			Comment("菜单规格项ID"),
		field.Uint64("category_spec_id").
			Optional().
			Nillable().
			Comment("菜单种类下的规格项ID"),
		// field.Uint64("")
		field.Float("price_delta").
			Default(0).
			Comment("特殊加价"),
		field.Int("sort").
			Default(0).
			Comment("顺序"),
	}
}

func (MenuSpec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("menu", Menu.Type).
			Ref("menu_specs").
			Field("menu_id").
			Unique().
			Required(), // 菜单
		edge.From("category_spec", MenuCategorySpec.Type).
			Ref("menu_specs").
			Field("category_spec_id").
			Unique(), // 菜单种类下的的规格组
		edge.To("spec_item", SpecItem.Type).
			Unique().
			Field("spec_item_id"), // 规格值
	}
}
