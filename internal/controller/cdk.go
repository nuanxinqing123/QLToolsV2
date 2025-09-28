package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
	"github.com/nuanxinqing123/QLToolsV2/internal/schema"
	"github.com/nuanxinqing123/QLToolsV2/internal/service"
)

type CDKController struct {
	cdkService *service.CDKService
}

// NewCDKController 创建CDKController实例
func NewCDKController() *CDKController {
	return &CDKController{
		cdkService: service.NewCDKService(),
	}
}

// CDKRouter CDK相关路由注册
func (ctrl *CDKController) CDKRouter(router *gin.RouterGroup) {
	router.GET("/list", ctrl.GetCDKList)                // 获取CDK列表
	router.POST("/create", ctrl.AddCDK)                 // 创建CDK
	router.POST("/create/batch", ctrl.AddCDKBatch)      // 批量创建CDK
	router.PUT("/update", ctrl.UpdateCDK)               // 更新CDK
	router.DELETE("/:id", ctrl.DeleteCDK)               // 删除CDK
	router.POST("/toggle-status", ctrl.ToggleCDKStatus) // 切换CDK状态
}

// AddCDK 添加CDK
// @Summary 添加CDK
// @Description 添加新的CDK卡密，需要提供密钥和可用次数
// @Tags CDK管理
// @Accept json
// @Produce json
// @Param request body schema.AddCDKRequest true "添加CDK请求参数"
// @Success 200 {object} response.Data{data=schema.AddCDKResponse} "添加成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "添加失败"
// @Router /api/cdk/create [post]
// @Security ApiKeyAuth
func (ctrl *CDKController) AddCDK(c *gin.Context) {
	// 解析请求参数
	var req schema.AddCDKRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层添加CDK
	resp, err := ctrl.cdkService.AddCDK(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// AddCDKBatch 批量添加CDK
// @Summary 批量添加CDK
// @Description 批量生成CDK卡密
// @Tags CDK管理
// @Accept json
// @Produce json
// @Param request body schema.AddCDKBatchRequest true "批量添加CDK请求参数"
// @Success 200 {object} response.Data{data=schema.AddCDKBatchResponse} "添加成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "添加失败"
// @Router /api/cdk/create/batch [post]
// @Security ApiKeyAuth
func (ctrl *CDKController) AddCDKBatch(c *gin.Context) {
	// 解析请求参数
	var req schema.AddCDKBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层批量添加CDK
	resp, err := ctrl.cdkService.AddCDKBatch(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// UpdateCDK 更新CDK
// @Summary 更新CDK
// @Description 更新CDK卡密信息
// @Tags CDK管理
// @Accept json
// @Produce json
// @Param request body schema.UpdateCDKRequest true "更新CDK请求参数"
// @Success 200 {object} response.Data{data=schema.UpdateCDKResponse} "更新成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "CDK不存在"
// @Failure 500 {object} response.Data "更新失败"
// @Router /api/cdk/update [put]
// @Security ApiKeyAuth
func (ctrl *CDKController) UpdateCDK(c *gin.Context) {
	// 解析请求参数
	var req schema.UpdateCDKRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层更新CDK
	resp, err := ctrl.cdkService.UpdateCDK(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// GetCDKList 获取CDK列表
// @Summary 获取CDK列表
// @Description 分页获取CDK列表，支持按密钥搜索和状态筛选
// @Tags CDK管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param key query string false "CDK密钥（模糊搜索）"
// @Param is_enable query bool false "是否启用"
// @Success 200 {object} response.Data{data=schema.GetCDKListResponse} "获取成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 500 {object} response.Data "获取失败"
// @Router /api/cdk/list [get]
// @Security ApiKeyAuth
func (ctrl *CDKController) GetCDKList(c *gin.Context) {
	// 解析查询参数
	var req schema.GetCDKListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层获取CDK列表
	resp, err := ctrl.cdkService.GetCDKList(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// DeleteCDK 删除CDK
// @Summary 删除CDK
// @Description 根据CDK ID删除CDK（软删除）
// @Tags CDK管理
// @Accept json
// @Produce json
// @Param id path int true "CDK ID"
// @Success 200 {object} response.Data{data=schema.DeleteCDKResponse} "删除成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "CDK不存在"
// @Failure 500 {object} response.Data "删除失败"
// @Router /api/cdk/{id} [delete]
// @Security ApiKeyAuth
func (ctrl *CDKController) DeleteCDK(c *gin.Context) {
	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "CDK ID格式错误")
		return
	}

	// 调用服务层删除CDK
	resp, err := ctrl.cdkService.DeleteCDK(schema.DeleteCDKRequest{ID: id})
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}

// ToggleCDKStatus 切换CDK启用状态
// @Summary 切换CDK状态
// @Description 启用或禁用CDK
// @Tags CDK管理
// @Accept json
// @Produce json
// @Param request body schema.ToggleCDKStatusRequest true "切换状态请求参数"
// @Success 200 {object} response.Data{data=schema.ToggleCDKStatusResponse} "切换成功"
// @Failure 400 {object} response.Data "请求参数错误"
// @Failure 404 {object} response.Data "CDK不存在"
// @Failure 500 {object} response.Data "切换失败"
// @Router /api/cdk/toggle-status [post]
// @Security ApiKeyAuth
func (ctrl *CDKController) ToggleCDKStatus(c *gin.Context) {
	// 解析请求参数
	var req schema.ToggleCDKStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResErrorWithMsg(c, response.CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层切换CDK状态
	resp, err := ctrl.cdkService.ToggleCDKStatus(req)
	if err != nil {
		response.ResErrorWithMsg(c, response.CodeGenericError, err.Error())
		return
	}

	response.ResSuccess(c, resp)
}
