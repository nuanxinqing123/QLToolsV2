package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/nuanxinqing123/QLToolsV2/internal/service"
)

type OpenController struct {
	service *service.OpenService
}

// NewOpenController 创建 OpenController
func NewOpenController() *OpenController {
	return &OpenController{
		service: service.NewOpenService(),
	}
}

// OpenRouter 公开路由
func (c *OpenController) OpenRouter(r *gin.RouterGroup) {
	r.POST("/check-cdk", c.CheckCDK)                   // 卡密检查
	r.GET("/services", c.GetOnlineServices)            // 获取在线服务
	r.GET("/slots/:env_id", c.CalculateAvailableSlots) // 计算剩余位置
	r.POST("/submit", c.SubmitVariable)                // 提交变量
}

// CheckCDK 检查卡密
// @Summary 检查卡密
// @Description 检查卡密是否存在、是否禁用、使用次数是否足够
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param request body schema.CheckCDKRequest true "检查卡密请求参数"
// @Success 200 {object} response.Data{data=schema.CheckCDKResponse} "检查成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "检查失败"
// @Router /open/check-cdk [post]
func (c *OpenController) CheckCDK(ctx *gin.Context) {
	// 解析请求参数
	var req schema.CheckCDKRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(ctx, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层检查卡密
	resp, err := c.service.CheckCDK(req)
	if err != nil {
		response.ResErrorWithMsg(ctx, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(ctx, resp)
}

// GetOnlineServices 获取在线服务
// @Summary 获取在线服务
// @Description 获取所有启用的环境变量数据，再查询变量下面所有绑定并启用的面板数据
// @Tags 公开接口
// @Accept json
// @Produce json
// @Success 200 {object} response.Data{data=schema.GetOnlineServicesResponse} "获取成功"
// @Failure 500 {object} response.Data "获取失败"
// @Router /open/services [get]
func (c *OpenController) GetOnlineServices(ctx *gin.Context) {
	// 直接调用服务层获取在线服务，不需要解析参数
	resp, err := c.service.GetOnlineServices()
	if err != nil {
		response.ResErrorWithMsg(ctx, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(ctx, resp)
}

// CalculateAvailableSlots 计算剩余位置
// @Summary 计算剩余位置
// @Description 根据变量配置的位置数，从绑定的面板中计算剩余可用位置数。如果位置数小于0，则固定返回0
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param env_id path int true "环境变量ID"
// @Success 200 {object} response.Data{data=schema.CalculateAvailableSlotsResponse} "计算成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "环境变量不存在"
// @Failure 500 {object} response.Data "计算失败"
// @Router /open/slots/{env_id} [get]
func (c *OpenController) CalculateAvailableSlots(ctx *gin.Context) {
	// 解析路径参数
	envIDStr := ctx.Param("env_id")
	envID, err := strconv.ParseInt(envIDStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(ctx, response.CodeInvalidParam, "环境变量ID格式错误")
		return
	}

	// 调用服务层计算剩余位置
	resp, err := c.service.CalculateAvailableSlots(schema.CalculateAvailableSlotsRequest{EnvID: envID})
	if err != nil {
		response.ResErrorWithMsg(ctx, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(ctx, resp)
}

// SubmitVariable 提交变量
// @Summary 提交变量
// @Description 提交环境变量，包含完整的验证流程：空内容检查、变量存在性检查、KEY验证、正则校验、位置计算、插件处理、数据提交
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param request body schema.SubmitVariableRequest true "提交变量请求参数"
// @Success 200 {object} response.Data{data=schema.SubmitVariableResponse} "提交成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "提交失败"
// @Router /open/submit [post]
func (c *OpenController) SubmitVariable(ctx *gin.Context) {
	// 解析请求参数
	var req schema.SubmitVariableRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(ctx, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层提交变量
	resp, err := c.service.SubmitVariable(req)
	if err != nil {
		response.ResErrorWithMsg(ctx, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(ctx, resp)
}
