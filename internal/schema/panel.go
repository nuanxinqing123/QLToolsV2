package schema

// AddPanelRequest 添加面板请求结构
type AddPanelRequest struct {
	Name         string `json:"name" binding:"required"`          // 面板名称
	URL          string `json:"url" binding:"required"`           // 连接地址
	ClientID     string `json:"client_id" binding:"required"`     // Client_ID
	ClientSecret string `json:"client_secret" binding:"required"` // Client_Secret
}

// AddPanelResponse 添加面板响应结构
type AddPanelResponse struct {
	ID      int64  `json:"id"`      // 面板ID
	Message string `json:"message"` // 消息
}

// UpdatePanelRequest 更新面板请求结构
type UpdatePanelRequest struct {
	ID           int64  `json:"id" binding:"required"`            // 面板ID
	Name         string `json:"name" binding:"required"`          // 面板名称
	URL          string `json:"url" binding:"required"`           // 连接地址
	ClientID     string `json:"client_id" binding:"required"`     // Client_ID
	ClientSecret string `json:"client_secret" binding:"required"` // Client_Secret
	IsEnable     *bool  `json:"is_enable"`                        // 是否启用（可选）
}

// UpdatePanelResponse 更新面板响应结构
type UpdatePanelResponse struct {
	Message string `json:"message"` // 消息
}

// GetPanelResponse 获取面板响应结构
type GetPanelResponse struct {
	ID           int64  `json:"id"`            // 面板ID
	Name         string `json:"name"`          // 面板名称
	URL          string `json:"url"`           // 连接地址
	ClientID     string `json:"client_id"`     // Client_ID
	ClientSecret string `json:"client_secret"` // Client_Secret（敏感信息，可考虑脱敏）
	IsEnable     bool   `json:"is_enable"`     // 是否启用
	Token        string `json:"token"`         // Token
	Params       int32  `json:"params"`        // Params
	CreatedAt    string `json:"created_at"`    // 创建时间
	UpdatedAt    string `json:"updated_at"`    // 更新时间
}

// GetPanelListRequest 获取面板列表请求结构
type GetPanelListRequest struct {
	Page     int    `form:"page" binding:"min=1"`              // 页码
	PageSize int    `form:"page_size" binding:"min=1,max=100"` // 每页数量
	Name     string `form:"name"`                              // 面板名称（模糊搜索）
	IsEnable *bool  `form:"is_enable"`                         // 是否启用
}

// GetPanelListResponse 获取面板列表响应结构
type GetPanelListResponse struct {
	Total int64              `json:"total"` // 总数
	List  []GetPanelResponse `json:"list"`  // 面板列表
}

// DeletePanelRequest 删除面板请求结构
type DeletePanelRequest struct {
	ID int64 `json:"id" binding:"required"` // 面板ID
}

// DeletePanelResponse 删除面板响应结构
type DeletePanelResponse struct {
	Message string `json:"message"` // 消息
}

// TogglePanelStatusRequest 切换面板状态请求结构
type TogglePanelStatusRequest struct {
	ID       int64 `json:"id" binding:"required"`        // 面板ID
	IsEnable bool  `json:"is_enable" binding:"required"` // 是否启用
}

// TogglePanelStatusResponse 切换面板状态响应结构
type TogglePanelStatusResponse struct {
	Message string `json:"message"` // 消息
}

// RefreshPanelTokenRequest 刷新面板Token请求结构
type RefreshPanelTokenRequest struct {
	ID int64 `json:"id" binding:"required"` // 面板ID
}

// RefreshPanelTokenResponse 刷新面板Token响应结构
type RefreshPanelTokenResponse struct {
	Message string `json:"message"` // 消息
	Token   string `json:"token"`   // 新Token
}

// TestPanelConnectionRequest 测试面板连接请求结构
type TestPanelConnectionRequest struct {
	URL          string `json:"url" binding:"required"`           // 连接地址
	ClientID     string `json:"client_id" binding:"required"`     // Client_ID
	ClientSecret string `json:"client_secret" binding:"required"` // Client_Secret
}

// TestPanelConnectionResponse 测试面板连接响应结构
type TestPanelConnectionResponse struct {
	Success     bool   `json:"success"`      // 连接是否成功
	Message     string `json:"message"`      // 消息
	Token       string `json:"token"`        // Token（连接成功时返回）
	Expiration  int    `json:"expiration"`   // Token过期时间（连接成功时返回）
	ResponseMsg string `json:"response_msg"` // API响应消息（连接失败时的详细信息）
}
