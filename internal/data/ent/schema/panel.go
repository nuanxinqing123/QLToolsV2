package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Panel 面板数据
type Panel struct {
	ent.Schema
}

// Fields of the Panel.
func (Panel) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable().Comment("主键ID"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.String("name").NotEmpty().Comment("名称"),
		field.String("url").NotEmpty().Comment("连接地址"),
		field.String("client_id").NotEmpty().Comment("Client_ID"),
		field.String("client_secret").NotEmpty().Comment("Client_Secret"),
		field.Bool("is_enable").Comment("是否启用"),
		field.String("token").NotEmpty().Comment("Token"),
		field.Int32("params").Comment("Params"),
	}
}

// Edges of the Panel.
func (Panel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("envs", Env.Type).
			Ref("panels"),
	}
}
