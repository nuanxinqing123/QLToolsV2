package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/nuanxinqing123/QLToolsV2/internal/service"
)

type EnvController struct {
	envService *service.EnvService
}

// NewEnvController 创建EnvController实例
func NewEnvController() *EnvController {
	return &EnvController{
		envService: service.NewEnvService(),
	}
}

// EnvRouter 变量相关路由注册
func (ctrl *EnvController) EnvRouter(router *gin.RouterGroup) {
	router.GET("/list", ctrl.GetEnvList)                // 获取变量列表
	router.GET("/:id", ctrl.GetEnv)                     // 获取单个变量信息
	router.POST("/create", ctrl.AddEnv)                 // 创建变量
	router.PUT("/update", ctrl.UpdateEnv)               // 更新变量
	router.DELETE("/:id", ctrl.DeleteEnv)               // 删除变量
	router.POST("/toggle-status", ctrl.ToggleEnvStatus) // 切换变量状态
	router.POST("/panels", ctrl.UpdateEnvPanels)        // 更新环境变量的面板绑定关系
	router.GET("/panels/:env_id", ctrl.GetEnvPanels)    // 获取变量关联的面板
}

// AddEnv 添加环境变量
// @Summary 添加环境变量
// @Description 添加新的环境变量配置
// @Tags 环境变量管理
// @Accept json
// @Produce json
// @Param request body schema.AddEnvRequest true "添加环境变量请求参数"
// @Success 200 {object} response.Data{data=schema.AddEnvResponse} "添加成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "添加失败"
// @Router /api/env/create [post]
// @Security ApiKeyAuth
func (ctrl *EnvController) AddEnv(c *gin.Context) {
	// 解析请求参数
	var req schema.AddEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层添加环境变量
	resp, err := ctrl.envService.AddEnv(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// UpdateEnv 更新环境变量
// @Summary 更新环境变量
// @Description 更新环境变量配置信息
// @Tags 环境变量管理
// @Accept json
// @Produce json
// @Param request body schema.UpdateEnvRequest true "更新环境变量请求参数"
// @Success 200 {object} response.Data{data=schema.UpdateEnvResponse} "更新成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "环境变量不存在"
// @Failure 500 {object} response.Data "更新失败"
// @Router /api/env/update [put]
// @Security ApiKeyAuth
func (ctrl *EnvController) UpdateEnv(c *gin.Context) {
	// 解析请求参数
	var req schema.UpdateEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层更新环境变量
	resp, err := ctrl.envService.UpdateEnv(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetEnv 获取单个环境变量信息
// @Summary 获取环境变量信息
// @Description 根据环境变量ID获取详细信息
// @Tags 环境变量管理
// @Accept json
// @Produce json
// @Param id path int true "环境变量ID"
// @Success 200 {object} response.Data{data=schema.GetEnvResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "环境变量不存在"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/env/{id} [get]
// @Security ApiKeyAuth
func (ctrl *EnvController) GetEnv(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "环境变量ID格式错误")
		return
	}

	// 调用服务层获取环境变量信息
	resp, err := ctrl.envService.GetEnv(id)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetEnvList 获取环境变量列表
// @Summary 获取环境变量列表
// @Description 分页获取环境变量列表，支持按名称搜索和状态筛选
// @Tags 环境变量管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "变量名称（模糊搜索）"
// @Param is_enable query bool false "是否启用"
// @Param mode query int false "模式筛选"
// @Success 200 {object} response.Data{data=schema.GetEnvListResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/env/list [get]
// @Security ApiKeyAuth
func (ctrl *EnvController) GetEnvList(c *gin.Context) {
	// 解析查询参数
	var req schema.GetEnvListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层获取环境变量列表
	resp, err := ctrl.envService.GetEnvList(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// DeleteEnv 删除环境变量
// @Summary 删除环境变量
// @Description 根据环境变量ID删除环境变量（软删除）
// @Tags 环境变量管理
// @Accept json
// @Produce json
// @Param id path int true "环境变量ID"
// @Success 200 {object} response.Data{data=schema.DeleteEnvResponse} "删除成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "环境变量不存在"
// @Failure 500 {object} response.Data "删除失败"
// @Router /api/env/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *EnvController) DeleteEnv(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "环境变量ID格式错误")
		return
	}

	// 调用服务层删除环境变量
	resp, err := ctrl.envService.DeleteEnv(schema.DeleteEnvConfigRequest{ID: id})
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// ToggleEnvStatus 切换环境变量启用状态
// @Summary 切换环境变量状态
// @Description 启用或禁用环境变量
// @Tags 环境变量管理
// @Accept json
// @Produce json
// @Param request body schema.ToggleEnvStatusRequest true "切换状态请求参数"
// @Success 200 {object} response.Data{data=schema.ToggleEnvStatusResponse} "切换成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "环境变量不存在"
// @Failure 500 {object} response.Data "切换失败"
// @Router /api/env/toggle-status [post]
// @Security ApiKeyAuth
func (ctrl *EnvController) ToggleEnvStatus(c *gin.Context) {
	// 解析请求参数
	var req schema.ToggleEnvStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层切换环境变量状态
	resp, err := ctrl.envService.ToggleEnvStatus(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// UpdateEnvPanels 更新环境变量面板绑定关系
// @Summary 更新环境变量面板绑定关系
// @Description 更新环境变量与面板的绑定关系，传入面板ID列表，系统会自动处理绑定和解绑
// @Tags 环境变量管理
// @Accept json
// @Produce json
// @Param request body schema.UpdateEnvPanelsRequest true "更新绑定关系请求参数"
// @Success 200 {object} response.Data{data=schema.UpdateEnvPanelsResponse} "更新成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "环境变量或面板不存在"
// @Failure 500 {object} response.Data "更新失败"
// @Router /api/env/panels [post]
// @Security ApiKeyAuth
func (ctrl *EnvController) UpdateEnvPanels(c *gin.Context) {
	// 解析请求参数
	var req schema.UpdateEnvPanelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层更新环境变量面板绑定关系
	resp, err := ctrl.envService.UpdateEnvPanels(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetEnvPanels 获取环境变量关联的面板
// @Summary 获取环境变量关联的面板
// @Description 获取指定环境变量关联的所有面板ID
// @Tags 环境变量管理
// @Accept json
// @Produce json
// @Param env_id path int true "环境变量ID"
// @Success 200 {object} response.Data{data=schema.GetEnvPanelsResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "环境变量不存在"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/env/panels/{env_id} [get]
// @Security ApiKeyAuth
func (ctrl *EnvController) GetEnvPanels(c *gin.Context) {
	// 解析路径参数
	envIDStr := c.Param("env_id")
	envID, err := strconv.ParseInt(envIDStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "环境变量ID格式错误")
		return
	}

	// 调用服务层获取环境变量关联的面板
	resp, err := ctrl.envService.GetEnvPanels(schema.GetEnvPanelsRequest{EnvID: envID})
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}
