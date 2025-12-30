package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PluginExecutionLog 插件执行日志表
type PluginExecutionLog struct {
	ent.Schema
}

// Fields of the PluginExecutionLog.
func (PluginExecutionLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable().Comment("主键ID"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("创建时间"),
		field.Int64("plugin_id").Comment("插件ID"),
		field.Int64("env_id").Comment("环境变量ID"),
		field.String("execution_status").Comment("执行状态(success,error,timeout)"),
		field.Int32("execution_time").Comment("执行耗时(毫秒)"),
		field.Text("input_data").Optional().Nillable().Comment("输入数据"),
		field.Text("output_data").Optional().Nillable().Comment("输出数据"),
		field.Text("error_message").Optional().Nillable().Comment("错误信息"),
		field.Text("stack_trace").Optional().Nillable().Comment("错误堆栈"),
	}
}

// Indexes of the PluginExecutionLog.
func (PluginExecutionLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
		index.Fields("plugin_id"),
		index.Fields("env_id"),
		index.Fields("execution_status"),
	}
}

// Edges of the PluginExecutionLog.
func (PluginExecutionLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("plugin", Plugin.Type).
			Ref("execution_logs").
			Field("plugin_id").
			Unique().
			Required(),
	}
}
