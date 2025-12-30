package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// EnvPlugin 环境变量插件关联表
type EnvPlugin struct {
	ent.Schema
}

// Fields of the EnvPlugin.
func (EnvPlugin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable().Comment("主键ID"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.Int64("env_id").Comment("环境变量ID"),
		field.Int64("plugin_id").Comment("插件ID"),
		field.Bool("is_enable").Default(true).Comment("是否启用"),
		field.Int32("execution_order").Default(100).Comment("执行顺序(数字越小越先执行)"),
		field.Text("config").Optional().Nillable().Comment("插件配置参数"),
	}
}

// Indexes of the EnvPlugin.
func (EnvPlugin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("env_id", "plugin_id").Unique(),
		index.Fields("env_id"),
		index.Fields("plugin_id"),
		index.Fields("is_enable"),
		index.Fields("execution_order"),
	}
}

// Edges of the EnvPlugin.
func (EnvPlugin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("env", Env.Type).
			Ref("env_plugins").
			Field("env_id").
			Unique().
			Required(),
		edge.From("plugin", Plugin.Type).
			Ref("env_plugins").
			Field("plugin_id").
			Unique().
			Required(),
	}
}
