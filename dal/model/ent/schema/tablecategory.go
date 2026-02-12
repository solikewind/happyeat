package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TableCategory 餐桌类别（如大厅、包间）
type TableCategory struct {
	ent.Schema
}

func (TableCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "table_categories"},
		schema.Comment("餐桌类别（如大厅、包间）"),
	}
}

func (TableCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(64).
			Comment("类别名称"),
		field.String("description").
			Optional().
			Nillable().
			Comment("描述"),
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

func (TableCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tables", Table.Type),
	}
}
