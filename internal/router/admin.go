package router

import (
	"github.com/gin-gonic/gin"

	"QLToolsV2/internal/controllers"
	"QLToolsV2/internal/middleware"
)

// InitRouterAdmin API
func InitRouterAdmin(r *gin.RouterGroup) {
	// 权限认证
	r.Use(middleware.Auth())

	// 用户路由
	user := controllers.UserController{}
	user.Router(r)

	// 变量路由
	env := controllers.EnvController{}
	env.Router(r)

	// 面板路由
	panel := controllers.PanelController{}
	panel.Router(r)

	// KEY路由
	key := controllers.KEYController{}
	key.Router(r)
}
