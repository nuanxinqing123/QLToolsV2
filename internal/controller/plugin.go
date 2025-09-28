package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/nuanxinqing123/QLToolsV2/internal/service"
)

type PluginController struct {
	pluginService *service.PluginService
}

// NewPluginController 创建PluginController实例
func NewPluginController() *PluginController {
	return &PluginController{
		pluginService: service.NewPluginService(),
	}
}

// PluginRouter 插件相关路由注册
func (ctrl *PluginController) PluginRouter(router *gin.RouterGroup) {
	router.GET("/list", ctrl.GetPluginList)                    // 获取插件列表
	router.GET("/:id", ctrl.GetPlugin)                         // 获取单个插件信息
	router.POST("/create", ctrl.CreatePlugin)                  // 创建插件
	router.PUT("/update", ctrl.UpdatePlugin)                   // 更新插件
	router.DELETE("/:id", ctrl.DeletePlugin)                   // 删除插件
	router.POST("/toggle-status", ctrl.TogglePluginStatus)     // 切换插件状态
	router.POST("/test", ctrl.TestPlugin)                      // 测试插件
	router.POST("/bind-env", ctrl.BindPluginToEnv)             // 绑定插件到环境变量
	router.POST("/unbind-env", ctrl.UnbindPluginFromEnv)       // 解绑插件与环境变量
	router.GET("/envs/:plugin_id", ctrl.GetPluginEnvs)         // 获取插件关联环境变量
	router.GET("/execution-logs", ctrl.GetPluginExecutionLogs) // 获取插件执行日志
}

// CreatePlugin 创建插件
// @Summary 创建插件
// @Description 创建新的插件，需要提供插件名称、脚本内容等信息
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param request body schema.CreatePluginRequest true "创建插件请求参数"
// @Success 200 {object} response.Data{data=schema.CreatePluginResponse} "创建成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "创建失败"
// @Router /api/plugin/create [post]
// @Security ApiKeyAuth
func (ctrl *PluginController) CreatePlugin(c *gin.Context) {
	// 解析请求参数
	var req schema.CreatePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层创建插件
	resp, err := ctrl.pluginService.CreatePlugin(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// UpdatePlugin 更新插件
// @Summary 更新插件
// @Description 更新插件配置信息
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param request body schema.UpdatePluginRequest true "更新插件请求参数"
// @Success 200 {object} response.Data{data=schema.UpdatePluginResponse} "更新成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "插件不存在"
// @Failure 500 {object} response.Data "更新失败"
// @Router /api/plugin/update [put]
// @Security ApiKeyAuth
func (ctrl *PluginController) UpdatePlugin(c *gin.Context) {
	// 解析请求参数
	var req schema.UpdatePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层更新插件
	resp, err := ctrl.pluginService.UpdatePlugin(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetPlugin 获取单个插件信息
// @Summary 获取插件信息
// @Description 根据插件ID获取插件详细信息
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param id path int true "插件ID"
// @Success 200 {object} response.Data{data=schema.GetPluginResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "插件不存在"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/plugin/{id} [get]
// @Security ApiKeyAuth
func (ctrl *PluginController) GetPlugin(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "插件ID格式错误")
		return
	}

	// 调用服务层获取插件信息
	resp, err := ctrl.pluginService.GetPlugin(id)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetPluginList 获取插件列表
// @Summary 获取插件列表
// @Description 分页获取插件列表，支持按名称搜索、触发事件和状态筛选
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "插件名称（模糊搜索）"
// @Param trigger_event query string false "触发事件"
// @Param is_enable query bool false "是否启用"
// @Success 200 {object} response.Data{data=schema.GetPluginListResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/plugin/list [get]
// @Security ApiKeyAuth
func (ctrl *PluginController) GetPluginList(c *gin.Context) {
	// 解析查询参数
	var req schema.GetPluginListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层获取插件列表
	resp, err := ctrl.pluginService.GetPluginList(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// DeletePlugin 删除插件
// @Summary 删除插件
// @Description 根据插件ID删除插件（软删除）
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param id path int true "插件ID"
// @Success 200 {object} response.Data{data=schema.DeletePluginResponse} "删除成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "插件不存在"
// @Failure 500 {object} response.Data "删除失败"
// @Router /api/plugin/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *PluginController) DeletePlugin(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "插件ID格式错误")
		return
	}

	// 调用服务层删除插件
	resp, err := ctrl.pluginService.DeletePlugin(schema.DeletePluginRequest{ID: id})
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// TogglePluginStatus 切换插件启用状态
// @Summary 切换插件状态
// @Description 启用或禁用插件
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param request body schema.TogglePluginStatusRequest true "切换状态请求参数"
// @Success 200 {object} response.Data{data=schema.TogglePluginStatusResponse} "切换成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "插件不存在"
// @Failure 500 {object} response.Data "切换失败"
// @Router /api/plugin/toggle-status [post]
// @Security ApiKeyAuth
func (ctrl *PluginController) TogglePluginStatus(c *gin.Context) {
	// 解析请求参数
	var req schema.TogglePluginStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层切换插件状态
	resp, err := ctrl.pluginService.TogglePluginStatus(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// TestPlugin 测试插件
// @Summary 测试插件
// @Description 测试插件脚本执行，验证脚本语法和逻辑
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param request body schema.TestPluginRequest true "测试插件请求参数"
// @Success 200 {object} response.Data{data=schema.TestPluginResponse} "测试完成"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "测试失败"
// @Router /api/plugin/test [post]
// @Security ApiKeyAuth
func (ctrl *PluginController) TestPlugin(c *gin.Context) {
	// 解析请求参数
	var req schema.TestPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层测试插件
	resp, err := ctrl.pluginService.TestPlugin(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// BindPluginToEnv 绑定插件到环境变量
// @Summary 绑定插件到环境变量
// @Description 将插件绑定到指定环境变量，可以配置插件参数和执行顺序
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param request body schema.BindPluginToEnvRequest true "绑定插件请求参数"
// @Success 200 {object} response.Data{data=schema.BindPluginToEnvResponse} "绑定成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "绑定失败"
// @Router /api/plugin/bind-env [post]
// @Security ApiKeyAuth
func (ctrl *PluginController) BindPluginToEnv(c *gin.Context) {
	// 解析请求参数
	var req schema.BindPluginToEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层绑定插件
	resp, err := ctrl.pluginService.BindPluginToEnv(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// UnbindPluginFromEnv 解绑插件与环境变量
// @Summary 解绑插件与环境变量
// @Description 解除插件与环境变量的绑定关系
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param request body schema.UnbindPluginFromEnvRequest true "解绑插件请求参数"
// @Success 200 {object} response.Data{data=schema.UnbindPluginFromEnvResponse} "解绑成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "解绑失败"
// @Router /api/plugin/unbind-env [post]
// @Security ApiKeyAuth
func (ctrl *PluginController) UnbindPluginFromEnv(c *gin.Context) {
	// 解析请求参数
	var req schema.UnbindPluginFromEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层解绑插件
	resp, err := ctrl.pluginService.UnbindPluginFromEnv(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetPluginEnvs 获取插件关联环境变量
// @Summary 获取插件关联环境变量
// @Description 获取指定插件关联的所有环境变量信息
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param plugin_id path int true "插件ID"
// @Success 200 {object} response.Data{data=schema.GetPluginEnvsResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/plugin/envs/{plugin_id} [get]
// @Security ApiKeyAuth
func (ctrl *PluginController) GetPluginEnvs(c *gin.Context) {
	// 解析路径参数
	pluginIDStr := c.Param("plugin_id")
	pluginID, err := strconv.ParseInt(pluginIDStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "插件ID格式错误")
		return
	}

	// 调用服务层获取插件关联环境变量
	resp, err := ctrl.pluginService.GetPluginEnvs(schema.GetPluginEnvsRequest{PluginID: pluginID})
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetPluginExecutionLogs 获取插件执行日志
// @Summary 获取插件执行日志
// @Description 分页获取插件执行日志，支持多条件筛选
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param plugin_id query int false "插件ID"
// @Param env_id query int false "环境变量ID"
// @Param execution_status query string false "执行状态"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Data{data=schema.GetPluginExecutionLogsResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/plugin/execution-logs [get]
// @Security ApiKeyAuth
func (ctrl *PluginController) GetPluginExecutionLogs(c *gin.Context) {
	// 解析查询参数
	var req schema.GetPluginExecutionLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层获取执行日志
	resp, err := ctrl.pluginService.GetPluginExecutionLogs(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}
