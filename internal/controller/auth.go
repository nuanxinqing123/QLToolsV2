package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/nuanxinqing123/QLToolsV2/internal/service"
)

// 使用内存存储验证码，适用于单实例场景；如为多实例部署请替换为共享存储（如 Redis）
var captchaStore = base64Captcha.DefaultMemStore

type AuthController struct {
	authService *service.AuthService
}

// NewAuthController 创建认证控制器实例
func NewAuthController() *AuthController {
	return &AuthController{
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
