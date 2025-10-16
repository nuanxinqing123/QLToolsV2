package controller

import (
	"image/color"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/nuanxinqing123/QLToolsV2/internal/service"
)

type AuthController struct {
	authService  *service.AuthService
	captchaStore base64Captcha.Store
}

type AuthRequiredController struct {
	authService *service.AuthService
}

// NewAuthController 创建认证控制器实例
func NewAuthController() *AuthController {
	return &AuthController{
		authService: service.NewAuthService(),
		// 使用内存存储验证码，适用于单实例场景；如为多实例部署请替换为共享存储（如 Redis）
		captchaStore: base64Captcha.DefaultMemStore,
	}
}

// NewAuthRequiredController 创建认证控制器实例
func NewAuthRequiredController() *AuthRequiredController {
	return &AuthRequiredController{
		authService: service.NewAuthService(),
	}
}

// AuthRouter 认证相关路由注册
func (ctrl *AuthController) AuthRouter(router *gin.RouterGroup) {
	// 无需认证的接口
	router.GET("/captcha", ctrl.GetCaptcha)    // 算术验证码
	router.POST("/login", ctrl.Login)          // 用户登录
	router.POST("/register", ctrl.Register)    // 用户注册
	router.POST("/refresh", ctrl.RefreshToken) // 刷新Token
}

// AuthRequiredRouter 认证相关路由注册（携带token）
func (ctrl *AuthRequiredController) AuthRequiredRouter(router *gin.RouterGroup) {
	router.GET("/logout", ctrl.Logout) // 用户登出
}

// GetCaptcha 生成算术验证码并返回ID与Base64图片
// @Summary 获取验证码
// @Description 生成算术验证码并返回验证码ID与Base64图片
// @Tags 认证管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Data{data=schema.GetCaptchaResponse} "验证码信息"
// @Failure 500 {object} response.Data "验证码生成失败"
// @Router /api/auth/captcha [get]
func (ctrl *AuthController) GetCaptcha(ctx *gin.Context) {
	// 算术验证码驱动配置
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverMath

	// NumOfDigits 表示算式的位数，例如 2 表示 a+b 结果在两位数内；可按需调整难度
	captchaConfig := base64Captcha.DriverMath{
		Height:          36,
		Width:           120,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	driverString = captchaConfig
	driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, ctrl.captchaStore)

	id, b64s, _, err := captcha.Generate()
	if err != nil {
		response.ResErrorWithMsg(ctx, response.CodeGenericError, "验证码生成失败")
		return
	}

	// 返回前端可直接显示的 Base64 图片
	response.ResSuccess(ctx, schema.GetCaptchaResponse{
		CaptchaID:     id,
		CaptchaBase64: b64s,
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口，需要验证码验证
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body schema.RegisterRequest true "注册请求参数"
// @Success 200 {object} response.Data{data=schema.RegisterResponse} "注册成功"
// @Failure 400 {object} response.Data "请求参数错误或验证码错误"
// @Failure 500 {object} response.Data "注册失败"
// @Router /api/auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	// 解析请求参数
	var req schema.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 验证验证码
	if !ctrl.captchaStore.Verify(req.CaptchaID, req.CaptchaCode, true) {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "验证码错误")
		return
	}

	// 调用服务层进行用户注册
	if err := ctrl.authService.Register(req); err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, schema.RegisterResponse{
		Message: "注册成功",
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，需要验证码验证，返回访问令牌和刷新令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body schema.LoginRequest true "登录请求参数"
// @Success 200 {object} response.Data{data=schema.LoginResponse} "登录成功"
// @Failure 400 {object} response.Data "请求参数错误或验证码错误"
// @Failure 401 {object} response.Data "认证失败"
// @Router /api/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	// 解析请求参数
	var req schema.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 验证验证码
	if !ctrl.captchaStore.Verify(req.CaptchaID, req.CaptchaCode, true) {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "验证码错误")
		return
	}

	// 调用服务层进行登录验证
	loginResp, err := ctrl.authService.Login(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeNeedLogin, err.Error())
		return
	}

	// 返回登录成功响应
	response.ResSuccess(c, schema.LoginResponse{
		Message:      "登录成功",
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口，注销当前用户的访问令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Data{data=schema.LogoutResponse} "登出成功"
// @Failure 401 {object} response.Data "需要登录"
// @Failure 500 {object} response.Data "登出失败"
// @Router /api/auth/logout [post]
// @Security ApiKeyAuth
func (ctrl *AuthRequiredController) Logout(c *gin.Context) {
	// 调用服务层注销Token
	if err := ctrl.authService.Logout(); err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, "登出失败: "+err.Error())
		return
	}

	response.ResSuccess(c, schema.LogoutResponse{
		Message: "登出成功",
	})
}

// RefreshToken 刷新访问Token
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body schema.RefreshTokenRequest true "刷新令牌请求参数"
// @Success 200 {object} response.Data{data=schema.RefreshTokenResponse} "刷新成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 401 {object} response.Data "刷新令牌无效"
// @Router /api/auth/refresh [post]
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	// 解析请求参数
	var req schema.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层刷新Token
	newAccessToken, err := ctrl.authService.RefreshToken(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidToken, err.Error())
		return
	}

	response.ResSuccess(c, newAccessToken)
}
