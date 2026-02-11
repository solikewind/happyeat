package schema

import (
	"time"

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

func (Menu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "menus"},
		schema.Comment("菜单项"),
	}
}

func (Menu) Fields() []ent.Field {
	return []ent.Field{
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
		field.Float("price").
			Comment("价格"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
	}
}

func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", MenuCategory.Type).Ref("menus").Unique().Required(),
		edge.To("specs", MenuSpec.Type),
	}
}
