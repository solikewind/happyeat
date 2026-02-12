package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OrderItem 订单明细（含菜单快照，便于历史不变）
type OrderItem struct {
	ent.Schema
}

func (OrderItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "order_items"},
		schema.Comment("订单明细"),
	}
}

func (OrderItem) Fields() []ent.Field {
	return []ent.Field{
		field.String("menu_name").
			MaxLen(128).
			Comment("菜品名称快照"),
		field.Int("quantity").
			Default(1).
			Comment("数量"),
		field.Float("unit_price").
			Comment("单价（含规格加价）"),
		field.Float("amount").
			Comment("小计金额"),
		field.String("spec_info").
			MaxLen(256).
			Optional().
			Nillable().
			Comment("规格描述快照，如 大份,中辣"),
		field.Int("sort").
			Default(0).
			Comment("排序"),
	}
}

func (OrderItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("items").Unique().Required(),
		edge.From("menu", Menu.Type).Ref("order_items").Unique().Comment("关联菜单，可选；删除菜单不影响历史订单"),
	}
}
