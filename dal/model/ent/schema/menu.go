package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Menu holds the schema definition for the Menu entity.
type Menu struct {
	ent.Schema
}

func (Menu) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Menu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "menus"},
		schema.Comment("菜单项"),
	}
}

func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("menu_category_id").
			Comment("菜单分类ID"),
		field.Uint64("object_id").
			Optional().
			Nillable().
			Comment("封面图对象ID（objects 表），与 image URL 二选一或同时存在"),
		field.String("name").
			MaxLen(128).
			Comment("菜名"),
		field.String("description").
			Optional().
			Nillable().
			Comment("描述"),
		field.String("image").
			MaxLen(512).
			Optional().
			Nillable().
			Comment("图片URL"),
		field.Int64("price").
			Comment("价格"),
	}
}

func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", MenuCategory.Type).
			Ref("menus").
			Field("menu_category_id").
			Unique().
			Required(),
		edge.From("cover_object", Object.Type).
			Ref("menu_covers").
			Field("object_id").
			Unique().
			Comment("封面图对象"),
		edge.To("menu_specs", MenuSpec.Type),
		edge.To("order_items", OrderItem.Type).
			Comment("被订单项引用，可选"),
	}
}
