package model

import (
	"gorm.io/gorm"
)

// CdKey 卡密数据
type CdKey struct {
	/*
		字段标签: https://gorm.io/zh_CN/docs/models.html
	*/
	gorm.Model
	// KEY值
	Key string `gorm:"column:key;type:varchar(255);not null;uniqueIndex;comment:KEY值" json:"key"`
	// 可用次数
	Count int `gorm:"column:count;type:int(11);not null;comment:可用次数" json:"count"`
	// 是否启用
	IsEnable bool `gorm:"column:is_enable;type:tinyint(1);default:1;comment:是否启用" json:"is_enable"`
}

type KeyCheck struct {
	Key string `json:"key" binding:"required"` // 卡密
}

// AddCDK 添加卡密
type AddCDK struct {
	Key   string `json:"key" binding:"required"`   // 卡密
	Count int    `json:"count" binding:"required"` // 使用次数
}

// BatchAddCDK 批量添加卡密
type BatchAddCDK struct {
	AddCount int    `json:"add_count" binding:"required"` // 生成数量
	Count    int    `json:"count" binding:"required"`     // 使用次数
	Prefix   string `json:"prefix"`                       // 前缀
}

// BatchOperationCDK 批量操作
type BatchOperationCDK struct {
	IDs      []int `json:"ids"`       // 卡密ID
	IsEnable bool  `json:"is_enable"` // 是否启用
}

// UpdateCDK 修改
type UpdateCDK struct {
	ID       int    `json:"id" binding:"required"` // ID
	Key      string `json:"key"`                   // 卡密
	Count    int    `json:"count"`                 // 使用次数
	IsEnable bool   `json:"is_enable"`             // 是否启用
}

// DeleteCDK 删除
type DeleteCDK struct {
	IDs []int `json:"ids"` // ID
}
