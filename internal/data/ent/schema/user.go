package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// User 用户数据
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			Comment("主键ID"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
		field.String("username").
			NotEmpty().
			Comment("用户名"),
		field.String("password").
			NotEmpty().
			Sensitive().
			Comment("密码"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
