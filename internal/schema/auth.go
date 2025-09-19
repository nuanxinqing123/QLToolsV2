package schema

// GetCaptchaResponse 获取验证码响应结构
type GetCaptchaResponse struct {
	CaptchaID     string `json:"captcha_id"`     // 验证码ID
	CaptchaBase64 string `json:"captcha_base64"` // 验证码图片Base64
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`     // 用户名
	Password    string `json:"password" binding:"required"`     // 密码
	CaptchaID   string `json:"captcha_id" binding:"required"`   // 验证码ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // 验证码
}

// RegisterResponse 注册响应结构
type RegisterResponse struct {
	Message string `json:"message"` // 消息
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username    string `json:"username" binding:"required"`     // 用户名
	Password    string `json:"password" binding:"required"`     // 密码
	CaptchaID   string `json:"captcha_id" binding:"required"`   // 验证码ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // 验证码
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Message      string `json:"message"`       // 消息
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
}

// LogoutResponse 登出响应结构
type LogoutResponse struct {
	Message string `json:"message"` // 消息
}

// RefreshTokenRequest 刷新Token请求结构
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // 刷新Token
}

// RefreshTokenResponse 刷新Token响应结构
type RefreshTokenResponse struct {
	Message     string `json:"message"`      // 消息
	AccessToken string `json:"access_token"` // 访问令牌
}

// GetProfileResponse 用户信息响应结构
type GetProfileResponse struct {
	UserID       int64  `json:"user_id"`       // 用户ID
	UserName     string `json:"username"`      // 用户名
	EnterpriseID int64  `json:"enterprise_id"` // 企业ID
	Role         string `json:"role"`          // 角色/身份
}
