package schema

import "encoding/json"

// CreatePluginRequest 创建插件请求结构
type CreatePluginRequest struct {
	Name             string `json:"name" binding:"required"`                       // 插件名称
	Description      string `json:"description"`                                   // 插件描述
	Version          string `json:"version" binding:"required"`                    // 插件版本
	Author           string `json:"author"`                                        // 插件作者
	ScriptContent    string `json:"script_content" binding:"required"`             // JavaScript脚本内容
	TriggerEvent     string `json:"trigger_event" binding:"required"`              // 触发事件
	ExecutionTimeout int    `json:"execution_timeout" binding:"min=100,max=30000"` // 执行超时时间(毫秒)
	Priority         int    `json:"priority" binding:"min=1,max=1000"`             // 执行优先级
}

// CreatePluginResponse 创建插件响应结构
type CreatePluginResponse struct {
	ID      int64  `json:"id"`      // 插件ID
	Message string `json:"message"` // 消息
}

// UpdatePluginRequest 更新插件请求结构
type UpdatePluginRequest struct {
	ID               int64  `json:"id" binding:"required"`                         // 插件ID
	Name             string `json:"name" binding:"required"`                       // 插件名称
	Description      string `json:"description"`                                   // 插件描述
	Version          string `json:"version" binding:"required"`                    // 插件版本
	Author           string `json:"author"`                                        // 插件作者
	ScriptContent    string `json:"script_content" binding:"required"`             // JavaScript脚本内容
	TriggerEvent     string `json:"trigger_event" binding:"required"`              // 触发事件
	ExecutionTimeout int    `json:"execution_timeout" binding:"min=100,max=30000"` // 执行超时时间(毫秒)
	Priority         int    `json:"priority" binding:"min=1,max=1000"`             // 执行优先级
	IsEnable         *bool  `json:"is_enable"`                                     // 是否启用（可选）
}

// UpdatePluginResponse 更新插件响应结构
type UpdatePluginResponse struct {
	Message string `json:"message"` // 消息
}

// GetPluginResponse 获取插件响应结构
type GetPluginResponse struct {
	ID               int64  `json:"id"`                // 插件ID
	Name             string `json:"name"`              // 插件名称
	Description      string `json:"description"`       // 插件描述
	Version          string `json:"version"`           // 插件版本
	Author           string `json:"author"`            // 插件作者
	ScriptContent    string `json:"script_content"`    // JavaScript脚本内容
	TriggerEvent     string `json:"trigger_event"`     // 触发事件
	IsEnable         bool   `json:"is_enable"`         // 是否启用
	ExecutionTimeout int    `json:"execution_timeout"` // 执行超时时间(毫秒)
	Priority         int    `json:"priority"`          // 执行优先级
	CreatedAt        string `json:"created_at"`        // 创建时间
	UpdatedAt        string `json:"updated_at"`        // 更新时间
}

// GetPluginListRequest 获取插件列表请求结构
type GetPluginListRequest struct {
	Page         int    `form:"page" binding:"min=1"`              // 页码
	PageSize     int    `form:"page_size" binding:"min=1,max=100"` // 每页数量
	Name         string `form:"name"`                              // 插件名称（模糊搜索）
	TriggerEvent string `form:"trigger_event"`                     // 触发事件
	IsEnable     *bool  `form:"is_enable"`                         // 是否启用
}

// GetPluginListResponse 获取插件列表响应结构
type GetPluginListResponse struct {
	Total int64               `json:"total"` // 总数
	List  []GetPluginResponse `json:"list"`  // 插件列表
}

// DeletePluginRequest 删除插件请求结构
type DeletePluginRequest struct {
	ID int64 `json:"id" binding:"required"` // 插件ID
}

// DeletePluginResponse 删除插件响应结构
type DeletePluginResponse struct {
	Message string `json:"message"` // 消息
}

// TogglePluginStatusRequest 切换插件状态请求结构
type TogglePluginStatusRequest struct {
	ID       int64 `json:"id" binding:"required"`        // 插件ID
	IsEnable bool  `json:"is_enable" binding:"required"` // 是否启用
}

// TogglePluginStatusResponse 切换插件状态响应结构
type TogglePluginStatusResponse struct {
	Message string `json:"message"` // 消息
}

// TestPluginRequest 测试插件请求结构
type TestPluginRequest struct {
	ScriptContent string `json:"script_content" binding:"required"` // JavaScript脚本内容
	TestEnvValue  string `json:"test_env_value"`                    // 测试环境变量值
}

// TestPluginResponse 测试插件响应结构
type TestPluginResponse struct {
	Success       bool            `json:"success"`        // 执行是否成功
	ExecutionTime int             `json:"execution_time"` // 执行耗时(毫秒)
	OutputData    json.RawMessage `json:"output_data"`    // 输出数据
	ErrorMessage  string          `json:"error_message"`  // 错误信息
}

// BindPluginToEnvRequest 绑定插件到环境变量请求结构
type BindPluginToEnvRequest struct {
	PluginID       int64           `json:"plugin_id" binding:"required"` // 插件ID
	EnvID          int64           `json:"env_id" binding:"required"`    // 环境变量ID
	ExecutionOrder int32           `json:"execution_order"`              // 执行顺序
	Config         json.RawMessage `json:"config"`                       // 插件配置参数
}

// BindPluginToEnvResponse 绑定插件到环境变量响应结构
type BindPluginToEnvResponse struct {
	Message string `json:"message"` // 消息
}

// UnbindPluginFromEnvRequest 解绑插件与环境变量请求结构
type UnbindPluginFromEnvRequest struct {
	PluginID int64 `json:"plugin_id" binding:"required"` // 插件ID
	EnvID    int64 `json:"env_id" binding:"required"`    // 环境变量ID
}

// UnbindPluginFromEnvResponse 解绑插件与环境变量响应结构
type UnbindPluginFromEnvResponse struct {
	Message string `json:"message"` // 消息
}

// GetPluginEnvsRequest 获取插件关联环境变量请求结构
type GetPluginEnvsRequest struct {
	PluginID int64 `form:"plugin_id" binding:"required"` // 插件ID
}

// GetPluginEnvsResponse 获取插件关联环境变量响应结构
type GetPluginEnvsResponse struct {
	PluginID int64                   `json:"plugin_id"` // 插件ID
	Envs     []PluginEnvRelationInfo `json:"envs"`      // 关联环境变量列表
}

// PluginEnvRelationInfo 插件环境变量关联信息
type PluginEnvRelationInfo struct {
	EnvID          int64           `json:"env_id"`          // 环境变量ID
	EnvName        string          `json:"env_name"`        // 环境变量名称
	IsEnable       bool            `json:"is_enable"`       // 是否启用
	ExecutionOrder int32           `json:"execution_order"` // 执行顺序
	Config         json.RawMessage `json:"config"`          // 插件配置参数
	CreatedAt      string          `json:"created_at"`      // 创建时间
}

// GetPluginExecutionLogsRequest 获取插件执行日志请求结构
type GetPluginExecutionLogsRequest struct {
	Page            int    `form:"page" binding:"min=1"`              // 页码
	PageSize        int    `form:"page_size" binding:"min=1,max=100"` // 每页数量
	PluginID        *int64 `form:"plugin_id"`                         // 插件ID
	EnvID           *int64 `form:"env_id"`                            // 环境变量ID
	ExecutionStatus string `form:"execution_status"`                  // 执行状态
	StartTime       string `form:"start_time"`                        // 开始时间
	EndTime         string `form:"end_time"`                          // 结束时间
}

// GetPluginExecutionLogsResponse 获取插件执行日志响应结构
type GetPluginExecutionLogsResponse struct {
	Total int64                    `json:"total"` // 总数
	List  []PluginExecutionLogInfo `json:"list"`  // 日志列表
}

// PluginExecutionLogInfo 插件执行日志信息
type PluginExecutionLogInfo struct {
	ID              int64           `json:"id"`               // 日志ID
	PluginID        int64           `json:"plugin_id"`        // 插件ID
	PluginName      string          `json:"plugin_name"`      // 插件名称
	EnvID           int64           `json:"env_id"`           // 环境变量ID
	EnvName         string          `json:"env_name"`         // 环境变量名称
	ExecutionStatus string          `json:"execution_status"` // 执行状态
	ExecutionTime   int             `json:"execution_time"`   // 执行耗时(毫秒)
	InputData       json.RawMessage `json:"input_data"`       // 输入数据
	OutputData      json.RawMessage `json:"output_data"`      // 输出数据
	ErrorMessage    string          `json:"error_message"`    // 错误信息
	CreatedAt       string          `json:"created_at"`       // 创建时间
}
