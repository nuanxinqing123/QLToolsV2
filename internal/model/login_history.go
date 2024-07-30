package model

import (
	"gorm.io/gorm"
)

// LoginHistory 登录历史
type LoginHistory struct {
	/*
		字段标签: https://gorm.io/zh_CN/docs/models.html
	*/
	gorm.Model
	// IP地址
	IP string `gorm:"column:ip;type:varchar(255);not null;comment:IP地址" json:"ip"`
	// 物理地址
	Address string `gorm:"column:address;type:varchar(255);comment:物理地址" json:"address"`
	// 状态	[0:失败 1:成功]
	State bool `gorm:"column:state;type:tinyint(1);default:0;comment:状态" json:"state"`
}
