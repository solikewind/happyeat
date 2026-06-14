package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/solikewind/happyeat/common/consts/enum"
)

// Settlement 结账单（多笔订单合并结账）
type Settlement struct {
	ent.Schema
}

func (Settlement) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Settlement) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "settlements"},
		schema.Comment("结账单"),
	}
}

func (Settlement) Fields() []ent.Field {
	return []ent.Field{
		field.String("customer_name").
			MaxLen(64).
			Comment("客户名"),
		field.Enum("status").
			GoType(enum.SettlementStatusUnsettled).
			Annotations(
				entsql.Default(string(enum.SettlementStatusUnsettled)),
			).
			Comment("unsettled=未结账 settled=已结账"),
		field.Int64("total_amount").
			Default(0).
			Comment("应收合计（分，关联订单 total_amount 之和）"),
		field.Int64("actual_amount").
			Default(0).
			Comment("实收合计（分，结账时录入）"),
		field.String("remark").
			MaxLen(512).
			Optional().
			Nillable().
			Comment("备注"),
		field.Time("settled_at").
			Optional().
			Nillable().
			Comment("结账时间"),
	}
}

func (Settlement) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("orders", Order.Type).
			Comment("关联订单"),
	}
}
