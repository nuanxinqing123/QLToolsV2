package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// LoginHistory 登录历史
type LoginHistory struct {
	ent.Schema
}

// Fields of the LoginHistory.
func (LoginHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable().Comment("主键ID"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.String("ip").NotEmpty().Comment("IP地址"),
		field.String("address").Optional().Nillable().Comment("物理地址"),
		field.Bool("state").Comment("状态 0:失败 1:成功"),
	}
}

// Edges of the LoginHistory.
func (LoginHistory) Edges() []ent.Edge {
	return nil
}
