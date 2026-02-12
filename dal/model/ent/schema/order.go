package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Order 订单（堂食关联餐桌，打包外带无桌台）
type Order struct {
	ent.Schema
}

func (Order) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "orders"},
		schema.Comment("订单"),
	}
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.String("order_no").
			MaxLen(64).
			Unique().
			Comment("订单号"),
		field.String("order_type").
			Default("dine_in").
			MaxLen(32).
			Comment("dine_in=堂食 takeaway=打包外带"),
		field.String("status").
			Default("created").
			MaxLen(32).
			Comment("created=待支付 paid=已支付 preparing=制作中 completed=已完成 cancelled=已取消"),
		field.Float("total_amount").
			Default(0).
			Comment("订单总金额"),
		field.String("remark").
			MaxLen(512).
			Optional().
			Nillable().
			Comment("备注"),
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

func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("table", Table.Type).Ref("orders").Unique().Comment("堂食时关联餐桌，外带为空"),
		edge.To("items", OrderItem.Type),
	}
}
