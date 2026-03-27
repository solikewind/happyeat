package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// IAMUser 权限系统用户主体。
type IAMUser struct {
	ent.Schema
}

func (IAMUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (IAMUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "iam_users"},
		schema.Comment("RBAC 用户主体"),
	}
}

func (IAMUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_code").
			NotEmpty().
			Unique().
			Immutable().
			Comment("用户编码"),
		field.String("display_name").
			Default("").
			Comment("展示名"),
	}
}

func (IAMUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", IAMRole.Type).
			StorageKey(
				edge.Table("iam_user_roles"),
				edge.Columns("iam_user_id", "iam_role_id"),
			),
	}
}
