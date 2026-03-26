package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/solikewind/happyeat/common/consts/enum"
)

// Order 订单（堂食关联餐桌，打包外带无桌台）
type Order struct {
	ent.Schema
}

func (Order) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Order) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "orders"},
		schema.Comment("订单表"),
	}
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("table_id").
			Optional().
			Nillable().
			Comment("餐桌ID"),
		field.String("order_no").
			MaxLen(64).
			Unique().
			Comment("订单号"),
		field.Enum("order_type").
			GoType(enum.OrderTypeDineIn).
			Annotations(
				entsql.Default(enum.OrderTypeDineIn.String()),
			).
			Comment("dine_in=堂食 takeaway=打包外带"),
		field.Enum("status").
			GoType(enum.OrderStatusCreated).
			Annotations(
				entsql.Default(enum.OrderStatusCreated.String()),
			).
			Comment("created=待支付 paid=已支付 preparing=制作中 completed=已完成 cancelled=已取消"),
		field.Int64("total_amount").
			Default(0).
			Comment("订单总金额"),
		field.String("remark").
			MaxLen(512).
			Optional().
			Nillable().
			Comment("备注"),
	}
}

func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("table", Table.Type).
			Ref("orders").
			Field("table_id").
			Unique().
			Comment("堂食时关联餐桌，外带为空"),
		edge.To("items", OrderItem.Type),
	}
}
