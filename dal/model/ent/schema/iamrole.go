package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// IAMRole 权限系统角色。
type IAMRole struct {
	ent.Schema
}

func (IAMRole) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (IAMRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "iam_roles"},
		schema.Comment("RBAC 角色"),
	}
}

func (IAMRole) Fields() []ent.Field {
	return []ent.Field{
		field.String("role_code").
			NotEmpty().
			Unique().
			Immutable().
			Comment("角色编码"),
		field.String("role_name").
			Default("").
			Comment("角色名称"),
	}
}

func (IAMRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", IAMUser.Type).
			Ref("roles"),
		edge.To("permissions", IAMPermission.Type).
			StorageKey(
				edge.Table("iam_role_permissions"),
				edge.Columns("iam_role_id", "iam_permission_id"),
			),
	}
}
