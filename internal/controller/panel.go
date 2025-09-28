package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/nuanxinqing123/QLToolsV2/internal/service"
)

type PanelController struct {
	panelService *service.PanelService
}

// NewPanelController 创建PanelController实例
func NewPanelController() *PanelController {
	return &PanelController{
		panelService: service.NewPanelService(),
	}
}

// PanelRouter 面板相关路由注册
func (ctrl *PanelController) PanelRouter(router *gin.RouterGroup) {
	router.GET("/list", ctrl.GetPanelList)                    // 获取面板列表
	router.GET("/:id", ctrl.GetPanel)                         // 获取单个面板信息
	router.POST("/create", ctrl.AddPanel)                     // 创建面板
	router.PUT("/update", ctrl.UpdatePanel)                   // 更新面板
	router.DELETE("/:id", ctrl.DeletePanel)                   // 删除面板
	router.POST("/toggle-status", ctrl.TogglePanelStatus)     // 切换面板状态
	router.POST("/refresh-token", ctrl.RefreshPanelToken)     // 刷新面板Token
	router.POST("/test-connection", ctrl.TestPanelConnection) // 测试面板连接
}

// AddPanel 添加面板
// @Summary 添加面板
// @Description 添加新的面板配置，需要提供面板名称、连接地址、Client_ID和Client_Secret
// @Tags 面板管理
// @Accept json
// @Produce json
// @Param request body schema.AddPanelRequest true "添加面板请求参数"
// @Success 200 {object} response.Data{data=schema.AddPanelResponse} "添加成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "添加失败"
// @Router /api/panel/create [post]
// @Security ApiKeyAuth
func (ctrl *PanelController) AddPanel(c *gin.Context) {
	// 解析请求参数
	var req schema.AddPanelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层添加面板
	resp, err := ctrl.panelService.AddPanel(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// UpdatePanel 更新面板
// @Summary 更新面板
// @Description 更新面板配置信息
// @Tags 面板管理
// @Accept json
// @Produce json
// @Param request body schema.UpdatePanelRequest true "更新面板请求参数"
// @Success 200 {object} response.Data{data=schema.UpdatePanelResponse} "更新成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "面板不存在"
// @Failure 500 {object} response.Data "更新失败"
// @Router /api/panel/update [put]
// @Security ApiKeyAuth
func (ctrl *PanelController) UpdatePanel(c *gin.Context) {
	// 解析请求参数
	var req schema.UpdatePanelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层更新面板
	resp, err := ctrl.panelService.UpdatePanel(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetPanel 获取单个面板信息
// @Summary 获取面板信息
// @Description 根据面板ID获取面板详细信息
// @Tags 面板管理
// @Accept json
// @Produce json
// @Param id path int true "面板ID"
// @Success 200 {object} response.Data{data=schema.GetPanelResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "面板不存在"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/panel/{id} [get]
// @Security ApiKeyAuth
func (ctrl *PanelController) GetPanel(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "面板ID格式错误")
		return
	}

	// 调用服务层获取面板信息
	resp, err := ctrl.panelService.GetPanel(id)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetPanelList 获取面板列表
// @Summary 获取面板列表
// @Description 分页获取面板列表，支持按名称搜索和状态筛选
// @Tags 面板管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "面板名称（模糊搜索）"
// @Param is_enable query bool false "是否启用"
// @Success 200 {object} response.Data{data=schema.GetPanelListResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/panel/list [get]
// @Security ApiKeyAuth
func (ctrl *PanelController) GetPanelList(c *gin.Context) {
	// 解析查询参数
	var req schema.GetPanelListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层获取面板列表
	resp, err := ctrl.panelService.GetPanelList(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// DeletePanel 删除面板
// @Summary 删除面板
// @Description 根据面板ID删除面板（软删除）
// @Tags 面板管理
// @Accept json
// @Produce json
// @Param id path int true "面板ID"
// @Success 200 {object} response.Data{data=schema.DeletePanelResponse} "删除成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "面板不存在"
// @Failure 500 {object} response.Data "删除失败"
// @Router /api/panel/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *PanelController) DeletePanel(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "面板ID格式错误")
		return
	}

	// 调用服务层删除面板
	resp, err := ctrl.panelService.DeletePanel(schema.DeletePanelRequest{ID: id})
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// TogglePanelStatus 切换面板启用状态
// @Summary 切换面板状态
// @Description 启用或禁用面板
// @Tags 面板管理
// @Accept json
// @Produce json
// @Param request body schema.TogglePanelStatusRequest true "切换状态请求参数"
// @Success 200 {object} response.Data{data=schema.TogglePanelStatusResponse} "切换成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "面板不存在"
// @Failure 500 {object} response.Data "切换失败"
// @Router /api/panel/toggle-status [post]
// @Security ApiKeyAuth
func (ctrl *PanelController) TogglePanelStatus(c *gin.Context) {
	// 解析请求参数
	var req schema.TogglePanelStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层切换面板状态
	resp, err := ctrl.panelService.TogglePanelStatus(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// RefreshPanelToken 刷新面板Token
// @Summary 刷新面板Token
// @Description 重新获取面板的访问Token
// @Tags 面板管理
// @Accept json
// @Produce json
// @Param request body schema.RefreshPanelTokenRequest true "刷新Token请求参数"
// @Success 200 {object} response.Data{data=schema.RefreshPanelTokenResponse} "刷新成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "面板不存在"
// @Failure 500 {object} response.Data "刷新失败"
// @Router /api/panel/refresh-token [post]
// @Security ApiKeyAuth
func (ctrl *PanelController) RefreshPanelToken(c *gin.Context) {
	// 解析请求参数
	var req schema.RefreshPanelTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层刷新面板Token
	resp, err := ctrl.panelService.RefreshPanelToken(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// TestPanelConnection 测试面板连接
// @Summary 测试面板连接
// @Description 测试面板连接是否正常，验证连接地址、Client_ID和Client_Secret是否有效
// @Tags 面板管理
// @Accept json
// @Produce json
// @Param request body schema.TestPanelConnectionRequest true "测试连接请求参数"
// @Success 200 {object} response.Data{data=schema.TestPanelConnectionResponse} "测试完成"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "测试失败"
// @Router /api/panel/test-connection [post]
// @Security ApiKeyAuth
func (ctrl *PanelController) TestPanelConnection(c *gin.Context) {
	// 解析请求参数
	var req schema.TestPanelConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层测试面板连接
	resp, err := ctrl.panelService.TestPanelConnection(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}
