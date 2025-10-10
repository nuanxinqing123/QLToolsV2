package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
	"github.com/nuanxinqing123/QLToolsV2/internal/service"
)

type DashboardController struct {
	dashboardService *service.DashboardService
}

// NewDashboardController 创建DashboardController实例
func NewDashboardController() *DashboardController {
	return &DashboardController{
		dashboardService: service.NewDashboardService(),
	}
}

// DashboardRouter 仪表盘相关路由注册
func (ctrl *DashboardController) DashboardRouter(router *gin.RouterGroup) {
	router.GET("/overview", ctrl.GetOverview)              // 获取数据总览
	router.GET("/submit-trend", ctrl.GetSubmitTrend)       // 获取提交趋势
	router.GET("/recent-activity", ctrl.GetRecentActivity) // 获取最近活动
	router.GET("/resource-usage", ctrl.GetResourceUsage)   // 获取资源使用情况
}

// GetOverview 获取数据总览
// @Summary 获取数据总览
// @Description 获取在线服务、总面板数、活跃CDK、今日提交等统计数据
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Success 200 {object} response.Data{data=schema.OverviewResponse} "获取成功"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/dashboard/overview [get]
// @Security ApiKeyAuth
func (ctrl *DashboardController) GetOverview(c *gin.Context) {
	data, err := ctrl.dashboardService.GetOverview()
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, data)
}

// GetSubmitTrend 获取提交趋势
// @Summary 获取提交趋势
// @Description 获取最近7天的提交趋势数据（当前为模拟数据）
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Success 200 {object} response.Data{data=schema.SubmitTrendResponse} "获取成功"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/dashboard/submit-trend [get]
// @Security ApiKeyAuth
func (ctrl *DashboardController) GetSubmitTrend(c *gin.Context) {
	data, err := ctrl.dashboardService.GetSubmitTrend()
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, data)
}

// GetRecentActivity 获取最近活动
// @Summary 获取最近活动
// @Description 获取最近的系统活动记录（当前为模拟数据）
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Success 200 {object} response.Data{data=schema.RecentActivityResponse} "获取成功"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/dashboard/recent-activity [get]
// @Security ApiKeyAuth
func (ctrl *DashboardController) GetRecentActivity(c *gin.Context) {
	data, err := ctrl.dashboardService.GetRecentActivity()
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, data)
}

// GetResourceUsage 获取资源使用情况
// @Summary 获取资源使用情况
// @Description 获取CPU、内存、磁盘的使用情况
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Success 200 {object} response.Data{data=schema.ResourceUsageResponse} "获取成功"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/dashboard/resource-usage [get]
// @Security ApiKeyAuth
func (ctrl *DashboardController) GetResourceUsage(c *gin.Context) {
	data, err := ctrl.dashboardService.GetResourceUsage()
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, data)
}
