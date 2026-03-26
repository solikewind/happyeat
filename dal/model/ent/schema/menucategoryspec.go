package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// MenuCategorySpec holds the schema definition for MenuCategory-level specs.
type MenuCategorySpec struct {
	ent.Schema
}

func (MenuCategorySpec) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (MenuCategorySpec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "category_specs"},
		schema.Comment("菜单分类的规格模板"),
	}
}

func (MenuCategorySpec) Fields() []ent.Field {
	return []ent.Field{
		// field.Uint64("menu_category_id").
		// 	Comment("菜单分类ID"),
		field.String("spec_type").
			MaxLen(64).
			Comment("规格类型，如辣度、容量"),
		field.String("spec_value").
			MaxLen(64).
			Comment("规格选项，如微辣、大杯"),
		field.Float("price_delta").
			Default(0).
			Comment("加价"),
		field.Int("sort").
			Default(0).
			Comment("排序"),
	}
}

func (MenuCategorySpec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", MenuCategory.Type).
			Ref("category_specs").
			// Field("menu_category_id").
			Unique().
			Required(),
		edge.To("menu_specs", MenuSpec.Type),
	}
}
