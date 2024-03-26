package router

import (
	"github.com/gin-gonic/gin"

	"QLToolsV2/internal/controllers"
)

// InitRouterAdmin API
func InitRouterAdmin(r *gin.RouterGroup) {
	user := controllers.UserController{}
	user.Router(r)
}
