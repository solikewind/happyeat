package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// IAMPermission 权限点定义。
type IAMPermission struct {
	ent.Schema
}

func (IAMPermission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (IAMPermission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "iam_permissions"},
		schema.Comment("RBAC 权限点"),
	}
}

func (IAMPermission) Fields() []ent.Field {
	return []ent.Field{
		field.String("permission_code").
			NotEmpty().
			Unique().
			Immutable().
			Comment("权限编码"),
		field.String("description").
			Default("").
			Comment("权限说明"),
	}
}

func (IAMPermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("roles", IAMRole.Type).
			Ref("permissions"),
	}
}
