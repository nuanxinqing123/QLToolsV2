package model

import (
	"gorm.io/gorm"
)

// User 用户数据表
type User struct {
	/*
		字段标签: https://gorm.io/zh_CN/docs/models.html

		UserID 用户ID
		UserName 昵称/用户名
		PassWord 密码
		Avatar 头像
		Balance 余额
		Role 用户角色   // user 用户、admin 管理员
		Status 用户状态 // active、inactive、suspend
	*/
	gorm.Model
	UserID   string  `gorm:"column:user_id;type:varchar(255);not null;uniqueIndex;comment:用户ID" json:"user_id"`
	UserName string  `gorm:"column:username;type:varchar(255);unique;comment:用户名" json:"username"`
	PassWord string  `gorm:"column:password;type:varchar(255);not null;comment:密码" json:"password"`
	Avatar   string  `gorm:"column:avatar;type:varchar(255);comment:头像" json:"avatar"`
	Balance  float64 `gorm:"column:balance;type:decimal(10,2);not null;default:0;comment:余额" json:"balance"`
	Role     string  `gorm:"column:role;type:varchar(10);not null;default:user;comment:角色" json:"role"`
	Status   string  `gorm:"column:status;type:varchar(10);not null;default:inactive;comment:状态" json:"status"`
}

const (
	// UserRole 普通用户
	UserRole string = "user"
	// AdminRole 管理员
	AdminRole string = "admin"

	// Active 激活用户
	Active string = "active"
	// Inactive 未激活用户
	Inactive string = "inactive"
	// Suspend 封禁用户
	Suspend string = "suspend"
)

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
