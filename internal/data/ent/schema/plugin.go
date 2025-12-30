package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Plugin 插件表
type Plugin struct {
	ent.Schema
}

// Fields of the Plugin.
func (Plugin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable().Comment("主键ID"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.String("name").NotEmpty().Unique().Comment("插件名称"),
		field.Text("description").Optional().Nillable().Comment("插件描述"),
		field.String("version").Default("1.0.0").Comment("插件版本"),
		field.String("author").Optional().Nillable().Comment("插件作者"),
		field.Text("script_content").NotEmpty().Comment("JavaScript脚本内容"),
		field.Bool("is_enable").Default(true).Comment("是否启用"),
		field.Int32("execution_timeout").Default(10000).Comment("执行超时时间(毫秒)"),
		field.String("trigger_event").Default("before_submit").Comment("触发事件"),
		field.Int32("priority").Default(10).Comment("执行优先级"),
	}
}

// Indexes of the Plugin.
func (Plugin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
		index.Fields("is_enable"),
	}
}

// Edges of the Plugin.
func (Plugin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("env_plugins", EnvPlugin.Type),
		edge.To("execution_logs", PluginExecutionLog.Type),
	}
}
