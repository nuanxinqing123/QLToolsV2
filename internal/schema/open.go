package schema

// CheckCDKRequest 检查卡密请求结构
type CheckCDKRequest struct {
	Key string `json:"key" binding:"required"` // CDK密钥
}

// CheckCDKResponse 检查卡密响应结构
type CheckCDKResponse struct {
	Valid         bool   `json:"valid"`          // 是否有效
	RemainingUses int32  `json:"remaining_uses"` // 剩余使用次数
	Message       string `json:"message"`        // 消息
}

// OnlineServiceInfo 在线服务信息
type OnlineServiceInfo struct {
	ID             int64   `json:"id"`              // 环境变量ID
	Name           string  `json:"name"`            // 变量名称
	Remarks        *string `json:"remarks"`         // 备注
	Quantity       int32   `json:"quantity"`        // 负载数量
	EnableKey      bool    `json:"enable_key"`      // 是否启用KEY
	CdkLimit       int32   `json:"cdk_limit"`       // 单次消耗卡密额度
	IsPrompt       bool    `json:"is_prompt"`       // 是否提示
	PromptLevel    *string `json:"prompt_level"`    // 提示等级
	PromptContent  *string `json:"prompt_content"`  // 提示内容
	AvailableSlots int32   `json:"available_slots"` // 可用位置数
}

// GetOnlineServicesResponse 获取在线服务响应结构
type GetOnlineServicesResponse struct {
	Total int64               `json:"total"` // 总数
	List  []OnlineServiceInfo `json:"list"`  // 在线服务列表
}

// CalculateAvailableSlotsRequest 计算剩余位置请求结构
type CalculateAvailableSlotsRequest struct {
	EnvID int64 `form:"env_id" binding:"required"` // 环境变量ID
}

// CalculateAvailableSlotsResponse 计算剩余位置响应结构
type CalculateAvailableSlotsResponse struct {
	EnvID          int64 `json:"env_id"`          // 环境变量ID
	TotalSlots     int32 `json:"total_slots"`     // 总位置数
	UsedSlots      int32 `json:"used_slots"`      // 已使用位置数
	AvailableSlots int32 `json:"available_slots"` // 可用位置数
}

// SubmitVariableRequest 提交变量请求结构
type SubmitVariableRequest struct {
	EnvID   int64  `json:"env_id" binding:"required"` // 环境变量ID
	Value   string `json:"value" binding:"required"`  // 变量值
	Key     string `json:"key"`                       // CDK密钥（如果启用KEY验证）
	Remarks string `json:"remarks"`                   // 备注
}

// SubmitVariableResponse 提交变量响应结构
type SubmitVariableResponse struct {
	Success      bool   `json:"success"`       // 是否成功
	Message      string `json:"message"`       // 消息
	SubmittedTo  int32  `json:"submitted_to"`  // 提交到的面板数量
	RemainingCDK int32  `json:"remaining_cdk"` // 剩余CDK次数（如果使用了CDK）
}
