package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Table 餐桌
type Table struct {
	ent.Schema
}

func (Table) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "tables"},
		schema.Comment("餐桌"),
	}
}

func (Table) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			MaxLen(32).
			Unique().
			Comment("桌号"),
		field.String("status").
			Default("idle").
			MaxLen(32).
			Comment("idle=空闲 using=使用中 reserved=预留 cleaning=清洁中"),
		field.Int("capacity").
			Default(4).
			Comment("可坐人数"),
		field.String("qr_code").
			MaxLen(256).
			Optional().
			Nillable().
			Comment("二维码"),
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

func (Table) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", TableCategory.Type).Ref("tables").Unique().Required(),
	}
}
