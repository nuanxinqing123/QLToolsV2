package model

import (
	"gorm.io/gorm"
)

// User 用户数据
type User struct {
	/*
		字段标签: https://gorm.io/zh_CN/docs/models.html
	*/
	gorm.Model
	// 用户ID
	UserID string `gorm:"column:user_id;type:varchar(255);not null;comment:用户ID" json:"user_id"`
	// 用户名
	UserName string `gorm:"column:username;type:varchar(255);not null;comment:用户名" json:"username"`
	// 密码
	PassWord string `gorm:"column:password;type:varchar(255);not null;comment:密码" json:"password"`
}

// Login 登录用户
type Login struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 注册
type Register struct {
	UserName   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 确认密码
}
