package router

import (
	"github.com/gin-gonic/gin"

	"QLToolsV2/internal/controllers"
)

// InitRouterOpen Api
func InitRouterOpen(r *gin.RouterGroup) {
	open := controllers.OpenController{}
	open.Router(r)
}
