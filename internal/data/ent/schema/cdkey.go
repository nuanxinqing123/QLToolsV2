package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CdKey 卡密数据
type CdKey struct {
	ent.Schema
}

// Fields of the CdKey.
func (CdKey) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Immutable().Comment("主键ID"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
		field.String("key").NotEmpty().Unique().Comment("KEY值"),
		field.Int32("count").Comment("可用次数"),
		field.Bool("is_enable").Default(true).Comment("是否启用"),
	}
}

// Indexes of the CdKey.
func (CdKey) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key").Unique(),
	}
}

// Edges of the CdKey.
func (CdKey) Edges() []ent.Edge {
	return nil
}
