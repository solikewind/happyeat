package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// MenuSpec holds the schema definition for the MenuSpec entity.
type MenuSpec struct {
	ent.Schema
}

func (MenuSpec) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "menu_specs"},
		schema.Comment("菜单规格（如口味、大小及加价）"),
	}
}

func (MenuSpec) Fields() []ent.Field {
	return []ent.Field{
		field.String("spec_type").
			MaxLen(64).
			Comment("规格类型（如口味、大小）"),
		field.String("spec_value").
			MaxLen(64).
			Comment("规格值（如中辣、大份）"),
		field.Float("price_delta").
			Default(0).
			Comment("加价"),
		field.Int("sort").
			Default(0).
			Comment("排序"),
	}
}

func (MenuSpec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("menu", Menu.Type).Ref("specs").Unique().Required(),
	}
}
