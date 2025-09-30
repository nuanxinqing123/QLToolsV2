package schema

// AddCDKRequest 添加CDK请求结构
type AddCDKRequest struct {
	Key   string `json:"key" binding:"required"`   // CDK密钥
	Count int32  `json:"count" binding:"required"` // 可用次数
}

// AddCDKResponse 添加CDK响应结构
type AddCDKResponse struct {
	ID      int64  `json:"id"`      // CDK ID
	Message string `json:"message"` // 消息
}

// AddCDKBatchRequest 批量添加CDK请求结构
type AddCDKBatchRequest struct {
	Count    int32 `json:"count" binding:"required"`     // 生成数量
	UseCount int32 `json:"use_count" binding:"required"` // 每个CDK的可用次数
}

// AddCDKBatchResponse 批量添加CDK响应结构
type AddCDKBatchResponse struct {
	Count   int32    `json:"count"`   // 成功生成数量
	Keys    []string `json:"keys"`    // 生成的CDK列表
	Message string   `json:"message"` // 消息
}

// UpdateCDKRequest 更新CDK请求结构
type UpdateCDKRequest struct {
	ID       int64  `json:"id" binding:"required"`    // CDK ID
	Key      string `json:"key" binding:"required"`   // CDK密钥
	Count    int32  `json:"count" binding:"required"` // 可用次数
	IsEnable *bool  `json:"is_enable"`                // 是否启用（可选）
}

// UpdateCDKResponse 更新CDK响应结构
type UpdateCDKResponse struct {
	Message string `json:"message"` // 消息
}

// GetCDKResponse 获取CDK响应结构
type GetCDKResponse struct {
	ID        int64  `json:"id"`         // CDK ID
	Key       string `json:"key"`        // CDK密钥
	Count     int32  `json:"count"`      // 可用次数
	IsEnable  bool   `json:"is_enable"`  // 是否启用
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// GetCDKListRequest 获取CDK列表请求结构
type GetCDKListRequest struct {
	Page     int    `form:"page" binding:"min=1"`              // 页码
	PageSize int    `form:"page_size" binding:"min=1,max=100"` // 每页数量
	Key      string `form:"key"`                               // CDK密钥（模糊搜索）
	IsEnable *bool  `form:"is_enable"`                         // 是否启用
}

// GetCDKListResponse 获取CDK列表响应结构
type GetCDKListResponse struct {
	Total int64            `json:"total"` // 总数
	List  []GetCDKResponse `json:"list"`  // CDK列表
}

// DeleteCDKRequest 删除CDK请求结构
type DeleteCDKRequest struct {
	ID int64 `json:"id" binding:"required"` // CDK ID
}

// DeleteCDKResponse 删除CDK响应结构
type DeleteCDKResponse struct {
	Message string `json:"message"` // 消息
}

// ToggleCDKStatusRequest 切换CDK状态请求结构
type ToggleCDKStatusRequest struct {
	ID       int64 `json:"id" binding:"required"` // CDK ID
	IsEnable bool  `json:"is_enable"`             // 是否启用
}

// ToggleCDKStatusResponse 切换CDK状态响应结构
type ToggleCDKStatusResponse struct {
	Message string `json:"message"` // 消息
}
