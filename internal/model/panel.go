package model

import (
	"gorm.io/gorm"
)

// Panel 面板数据
type Panel struct {
	/*
		字段标签: https://gorm.io/zh_CN/docs/models.html
	*/
	gorm.Model
	// 名称
	Name string `gorm:"column:name;type:varchar(255);not null;comment:名称" json:"name"`
	// 连接地址
	URL string `gorm:"column:url;type:varchar(255);not null;comment:连接地址" json:"url"`
	// Client_ID
	ClientID string `gorm:"column:client_id;type:varchar(255);not null;comment:Client_ID" json:"client_id"`
	// Client_Secret
	ClientSecret string `gorm:"column:client_secret;type:varchar(255);not null;comment:Client_Secret" json:"client_secret"`
	// 是否启用
	IsEnable bool `gorm:"column:is_enable;type:tinyint(1);default:0;comment:是否启用" json:"is_enable"`
	// Token
	Token string `gorm:"column:token;type:varchar(255);not null;comment:Token" json:"token"`
	// Params
	Params int `gorm:"column:params;type:int(11);not null;comment:Params" json:"params"`
}

// TestPanel 测试连接
type TestPanel struct {
	URL          string `json:"url" binding:"required"`           // 连接地址
	ClientID     string `json:"client_id" binding:"required"`     // Client_ID
	ClientSecret string `json:"client_secret" binding:"required"` // Client_Secret
}

// AddPanel 添加面板
type AddPanel struct {
	Name         string `json:"name" binding:"required"`          // 名称
	URL          string `json:"url" binding:"required"`           // 连接地址
	ClientID     string `json:"client_id" binding:"required"`     // Client_ID
	ClientSecret string `json:"client_secret" binding:"required"` // Client_Secret
	IsEnable     bool   `json:"is_enable"`                        // 是否启用
}

// BatchOperationPanel 批量操作
type BatchOperationPanel struct {
	IDs      []int `json:"ids"`       // 面板ID
	IsEnable bool  `json:"is_enable"` // 是否启用
}

// UpdatePanel 修改
type UpdatePanel struct {
	ID           int    `json:"id" binding:"required"`            // ID
	Name         string `json:"name" binding:"required"`          // 名称
	URL          string `json:"url" binding:"required"`           // 连接地址
	ClientID     string `json:"client_id" binding:"required"`     // Client_ID
	ClientSecret string `json:"client_secret" binding:"required"` // Client_Secret
	IsEnable     bool   `json:"is_enable"`                        // 是否启用
}

// DeletePanel 删除
type DeletePanel struct {
	IDs []int `json:"ids"` // ID
}
