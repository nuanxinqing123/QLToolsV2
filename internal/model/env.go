package model

import (
	"gorm.io/gorm"
)

// Env 变量数据
type Env struct {
	/*
		字段标签: https://gorm.io/zh_CN/docs/models.html
	*/
	gorm.Model
	// 名称
	Name string `gorm:"column:name;type:varchar(255);not null;comment:名称" json:"name"`
	// 备注
	Remarks string `gorm:"column:remarks;type:varchar(255);comment:备注" json:"remarks"`
	// 负载数量
	Quantity int `gorm:"column:quantity;type:int(11);not null;comment:负载数量" json:"quantity"`
	// 匹配正则
	Regex string `gorm:"column:regex;type:text;comment:匹配正则" json:"regex"`
	// 模式[1：新建模式、2：合并模式、3、更新模式]
	Mode int `gorm:"column:mode;type:int(11);not null;comment:模式" json:"mode"`
	// 分隔符[合并模式]
	Division string `gorm:"column:division;type:text;comment:分隔符[合并]" json:"division"`
	// 匹配正则[更新模式]
	RegexUpdate string `gorm:"column:regex_update;type:text;comment:匹配正则[更新]" json:"regex_update"`
	// 是否启用KEY
	EnableKey bool `gorm:"column:enable_key;type:tinyint(1);default:0;comment:是否启用KEY" json:"enable_key"`
	// 是否启用
	IsEnable bool `gorm:"column:is_enable;type:tinyint(1);default:0;comment:是否启用" json:"is_enable"`

	// 关联面板[Many To Many]
	Panels []Panel `gorm:"many2many:env_panels;" json:"panels"`
}

// AddEnv 添加变量
type AddEnv struct {
	Name        string `json:"name" binding:"required"`     // 名称
	Remarks     string `json:"remarks"`                     // 备注
	Quantity    int    `json:"quantity" binding:"required"` // 负载数量
	Regex       string `json:"regex"`                       // 匹配正则
	Mode        int    `json:"mode" binding:"required"`     // 模式
	Division    string `json:"division"`                    // 分隔符
	RegexUpdate string `json:"regex_update"`                // 匹配正则
	EnableKey   bool   `json:"enable_key"`                  // 是否启用KEY
	IsEnable    bool   `json:"is_enable"`                   // 是否启用变量
}

// BatchOperationEnv 批量操作
type BatchOperationEnv struct {
	IDs      []int `json:"ids"`       // 变量ID
	IsEnable bool  `json:"is_enable"` // 是否启用
}

// UpdateEnv 修改
type UpdateEnv struct {
	ID          int    `json:"id" binding:"required"`       // ID
	Name        string `json:"name" binding:"required"`     // 名称
	Remarks     string `json:"remarks"`                     // 备注
	Quantity    int    `json:"quantity" binding:"required"` // 负载数量
	Regex       string `json:"regex"`                       // 匹配正则
	Mode        int    `json:"mode" binding:"required"`     // 模式
	Division    string `json:"division"`                    // 分隔符
	RegexUpdate string `json:"regex_update"`                // 匹配正则
	EnableKey   bool   `json:"enable_key"`                  // 是否启用KEY
	IsEnable    bool   `json:"is_enable"`                   // 是否启用
}

// BindPanel 绑定面板
type BindPanel struct {
	EnvID    int   `json:"env_id" binding:"required"`    // 变量ID
	PanelIDs []int `json:"panel_ids" binding:"required"` // 面板ID
}

// DeleteEnv 删除
type DeleteEnv struct {
	IDs []int `json:"ids"` // ID
}
