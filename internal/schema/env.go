package schema

// AddEnvRequest 添加环境变量请求结构
type AddEnvRequest struct {
	Name            string  `json:"name" binding:"required"`      // 变量名称
	Remarks         *string `json:"remarks"`                      // 备注
	Quantity        int32   `json:"quantity" binding:"required"`  // 负载数量
	Regex           *string `json:"regex"`                        // 匹配正则
	Mode            int32   `json:"mode" binding:"required"`      // 模式
	RegexUpdate     *string `json:"regex_update"`                 // 匹配正则[更新]
	IsAutoEnvEnable bool    `json:"is_auto_env_enable"`           // 是否自动启用提交的变量
	EnableKey       bool    `json:"enable_key"`                   // 是否启用KEY
	CdkLimit        int32   `json:"cdk_limit" binding:"required"` // 单次消耗卡密额度
	IsPrompt        bool    `json:"is_prompt"`                    // 是否提示
	PromptLevel     *string `json:"prompt_level"`                 // 提示等级
	PromptContent   *string `json:"prompt_content"`               // 提示内容
}

// AddEnvResponse 添加环境变量响应结构
type AddEnvResponse struct {
	ID      int64  `json:"id"`      // 环境变量ID
	Message string `json:"message"` // 消息
}

// UpdateEnvRequest 更新环境变量请求结构
type UpdateEnvRequest struct {
	ID              int64   `json:"id" binding:"required"`        // 环境变量ID
	Name            string  `json:"name" binding:"required"`      // 变量名称
	Remarks         *string `json:"remarks"`                      // 备注
	Quantity        int32   `json:"quantity" binding:"required"`  // 负载数量
	Regex           *string `json:"regex"`                        // 匹配正则
	Mode            int32   `json:"mode" binding:"required"`      // 模式
	RegexUpdate     *string `json:"regex_update"`                 // 匹配正则[更新]
	IsAutoEnvEnable bool    `json:"is_auto_env_enable"`           // 是否自动启用提交的变量
	EnableKey       bool    `json:"enable_key"`                   // 是否启用KEY
	CdkLimit        int32   `json:"cdk_limit" binding:"required"` // 单次消耗卡密额度
	IsPrompt        bool    `json:"is_prompt"`                    // 是否提示
	PromptLevel     *string `json:"prompt_level"`                 // 提示等级
	PromptContent   *string `json:"prompt_content"`               // 提示内容
	IsEnable        *bool   `json:"is_enable"`                    // 是否启用（可选）
}

// UpdateEnvResponse 更新环境变量响应结构
type UpdateEnvResponse struct {
	Message string `json:"message"` // 消息
}

// GetEnvResponse 获取环境变量响应结构
type GetEnvResponse struct {
	ID              int64   `json:"id"`                 // 环境变量ID
	Name            string  `json:"name"`               // 变量名称
	Remarks         *string `json:"remarks"`            // 备注
	Quantity        int32   `json:"quantity"`           // 负载数量
	Regex           *string `json:"regex"`              // 匹配正则
	Mode            int32   `json:"mode"`               // 模式
	RegexUpdate     *string `json:"regex_update"`       // 匹配正则[更新]
	IsAutoEnvEnable bool    `json:"is_auto_env_enable"` // 是否自动启用提交的变量
	EnableKey       bool    `json:"enable_key"`         // 是否启用KEY
	CdkLimit        int32   `json:"cdk_limit"`          // 单次消耗卡密额度
	IsPrompt        bool    `json:"is_prompt"`          // 是否提示
	PromptLevel     *string `json:"prompt_level"`       // 提示等级
	PromptContent   *string `json:"prompt_content"`     // 提示内容
	IsEnable        bool    `json:"is_enable"`          // 是否启用
	CreatedAt       string  `json:"created_at"`         // 创建时间
	UpdatedAt       string  `json:"updated_at"`         // 更新时间
}

// GetEnvListRequest 获取环境变量列表请求结构
type GetEnvListRequest struct {
	Page     int    `form:"page" binding:"min=1"`              // 页码
	PageSize int    `form:"page_size" binding:"min=1,max=100"` // 每页数量
	Name     string `form:"name"`                              // 变量名称（模糊搜索）
	IsEnable *bool  `form:"is_enable"`                         // 是否启用
	Mode     *int32 `form:"mode"`                              // 模式筛选
}

// GetEnvListResponse 获取环境变量列表响应结构
type GetEnvListResponse struct {
	Total int64            `json:"total"` // 总数
	List  []GetEnvResponse `json:"list"`  // 环境变量列表
}

// DeleteEnvConfigRequest 删除环境变量请求结构
type DeleteEnvConfigRequest struct {
	ID int64 `json:"id" binding:"required"` // 环境变量ID
}

// DeleteEnvConfigResponse 删除环境变量响应结构
type DeleteEnvConfigResponse struct {
	Message string `json:"message"` // 消息
}

// ToggleEnvStatusRequest 切换环境变量状态请求结构
type ToggleEnvStatusRequest struct {
	ID       int64 `json:"id" binding:"required"` // 环境变量ID
	IsEnable bool  `json:"is_enable"`             // 是否启用
}

// ToggleEnvStatusResponse 切换环境变量状态响应结构
type ToggleEnvStatusResponse struct {
	Message string `json:"message"` // 消息
}

// UpdateEnvPanelsRequest 更新环境变量面板绑定关系请求结构
type UpdateEnvPanelsRequest struct {
	EnvID    int64   `json:"env_id" binding:"required"`    // 环境变量ID
	PanelIDs []int64 `json:"panel_ids" binding:"required"` // 面板ID列表（空数组表示解绑所有面板）
}

// UpdateEnvPanelsResponse 更新环境变量面板绑定关系响应结构
type UpdateEnvPanelsResponse struct {
	Message string `json:"message"` // 消息
}

// GetEnvPanelsRequest 获取环境变量关联面板请求结构
type GetEnvPanelsRequest struct {
	EnvID int64 `form:"env_id" binding:"required"` // 环境变量ID
}

// GetEnvPanelsResponse 获取环境变量关联面板响应结构
type GetEnvPanelsResponse struct {
	EnvID    int64   `json:"env_id"`    // 环境变量ID
	PanelIDs []int64 `json:"panel_ids"` // 关联的面板ID列表
}

// GetEnvPluginsRequest 获取环境变量关联插件请求结构
type GetEnvPluginsRequest struct {
	EnvID int64 `form:"env_id" binding:"required"` // 环境变量ID
}

// GetEnvPluginsResponse 获取环境变量关联插件响应结构
type GetEnvPluginsResponse struct {
	EnvID   int64                   `json:"env_id"` // 环境变量ID
	Plugins []EnvPluginRelationInfo `json:"plugins"` // 关联插件列表
}

// EnvPluginRelationInfo 环境变量插件关联信息
type EnvPluginRelationInfo struct {
	PluginID       int64  `json:"plugin_id"`       // 插件ID
	PluginName     string `json:"plugin_name"`     // 插件名称
	IsEnable       bool   `json:"is_enable"`       // 是否启用
	ExecutionOrder int32  `json:"execution_order"` // 执行顺序
	Config         string `json:"config"`          // 插件配置参数
	CreatedAt      string `json:"created_at"`      // 创建时间
}
