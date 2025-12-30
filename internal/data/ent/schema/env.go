package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Env 变量数据
type Env struct {
	ent.Schema
}

// Fields of the Env.
func (Env) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable().Comment("主键ID"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.String("name").NotEmpty().Comment("名称"),
		field.String("remarks").Optional().Nillable().Comment("备注"),
		field.Int32("quantity").Comment("负载数量"),
		field.Text("regex").Optional().Nillable().Comment("匹配正则"),
		field.Int32("mode").Comment("模式"),
		field.Text("regex_update").Optional().Nillable().Comment("匹配正则[更新]"),
		field.Bool("is_auto_env_enable").Default(true).Comment("是否自动启用提交的变量"),
		field.Bool("enable_key").Comment("是否启用KEY"),
		field.Int32("cdk_limit").Default(1).Comment("单次消耗卡密额度"),
		field.Bool("is_prompt").Comment("是否提示"),
		field.String("prompt_level").Optional().Nillable().Comment("提示等级"),
		field.Text("prompt_content").Optional().Nillable().Comment("提示内容"),
		field.Bool("is_enable").Comment("是否启用"),
	}
}

// Edges of the Env.
func (Env) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("panels", Panel.Type).
			StorageKey(edge.Table("env_panels"), edge.Columns("env_id", "panel_id")),
		edge.To("env_plugins", EnvPlugin.Type),
	}
}
