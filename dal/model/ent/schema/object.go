package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Object 对象元数据（文件存储于 COS，DB 保存元信息）
type Object struct {
	ent.Schema
}

func (Object) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UniqueID{},
		TimeMixin{},
		SoftDeleteMixin{},
	}
}

func (Object) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "objects"},
		schema.Comment("对象元数据表"),
	}
}

func (Object) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(255).
			Comment("原始文件名"),
		field.String("key").
			MaxLen(512).
			Unique().
			Comment("对象存储 key"),
		field.String("url").
			MaxLen(1024).
			Comment("访问 URL"),
		field.String("content_type").
			MaxLen(128).
			Optional().
			Nillable().
			Comment("内容类型"),
		field.String("suffix").
			MaxLen(16).
			Optional().
			Nillable().
			Comment("文件后缀"),
		field.Int64("size").
			Default(0).
			Comment("文件大小（字节）"),
		field.String("hash").
			MaxLen(64).
			Comment("内容哈希（murmur3）"),
	}
}

func (Object) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("menu_covers", Menu.Type),
	}
}

func (Object) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("hash"),
	}
}
