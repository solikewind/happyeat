package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CategorySpec holds the schema definition for category-level specs.
type CategorySpec struct {
	ent.Schema
}

func (CategorySpec) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (CategorySpec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "category_specs"},
		schema.Comment("菜单分类的规格模板"),
	}
}

func (CategorySpec) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("menu_category_id").
			Comment("菜单分类ID"),
		field.Uint64("spec_item_id").
			Optional().
			Nillable().
			Comment("引用的全局规格项ID"),
		field.String("spec_type").
			MaxLen(64).
			Comment("规格类型，如辣度、容量"),
		field.String("spec_value").
			MaxLen(64).
			Comment("规格选项，如微辣、大杯"),
		field.Int64("price_delta").
			Default(0).
			Comment("加价"),
		field.Uint32("sort").
			Default(0).
			Comment("排序"),
	}
}

func (CategorySpec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", MenuCategory.Type).
			Ref("category_specs").
			Field("menu_category_id").
			Unique().
			Required(), // 菜单分类
		edge.From("spec_item", SpecItem.Type).
			Ref("category_specs").
			Field("spec_item_id").
			Unique(),
		edge.To("menu_specs", MenuSpec.Type), // 菜单规格
	}
}
