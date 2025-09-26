package controller

import (
	"github.com/gin-gonic/gin"
	res "github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
	"github.com/nuanxinqing123/QLToolsV2/internal/service"
)

type HealthyController struct {
	service *service.HealthyService
}

// NewHealthyController 创建 HealthyController
func NewHealthyController() *HealthyController {
	return &HealthyController{
		service: service.NewHealthyService(),
	}
}

// HealthyRouter 注册路由
func (c *HealthyController) HealthyRouter(r *gin.RouterGroup) {
	// 健康检查
	r.GET("/healthy", c.Healthy)
}

// Healthy 健康检查
func (c *HealthyController) Healthy(ctx *gin.Context) {
	// 业务处理
	resCode, msg := c.service.CheckHealth()
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}
