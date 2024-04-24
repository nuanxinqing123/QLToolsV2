package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"QLToolsV2/internal/model"
	"QLToolsV2/internal/service"
	res "QLToolsV2/pkg/response"
	val "QLToolsV2/utils/validator"
)

type EnvController struct{}

// Router 注册路由
func (c *EnvController) Router(r *gin.RouterGroup) {
	// 分页查询
	r.GET("/env/list", c.List)

	// 添加
	r.POST("/env/add", c.Add)
	// 批量操作[启用/禁用]
	r.PUT("/env/batch/operation", c.BatchOperation)
	// 修改
	r.PUT("/env/update", c.Update)
	// 绑定面板
	r.PUT("/env/bind/panel", c.BindPanel)
	// 删除
	r.DELETE("/env/delete", c.Delete)
}

// List 分页查询
func (c *EnvController) List(ctx *gin.Context) {
	// 获取参数
	p := new(model.Pagination)
	if err := ctx.ShouldBindQuery(&p); err != nil {
		// 判断err是不是validator.ValidationErrors类型
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			res.ResError(ctx, res.CodeInvalidParam)
			return
		}

		// 翻译错误
		res.ResErrorWithMsg(ctx, res.CodeInvalidParam, val.RemoveTopStruct(errs.Translate(val.Trans)))
		return
	}

	// 业务处理
	resCode, msg := service.EnvList(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}

// Add 添加
func (c *EnvController) Add(ctx *gin.Context) {
	// 获取参数
	p := new(model.AddEnv)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		// 判断err是不是validator.ValidationErrors类型
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			res.ResError(ctx, res.CodeInvalidParam)
			return
		}

		// 翻译错误
		res.ResErrorWithMsg(ctx, res.CodeInvalidParam, val.RemoveTopStruct(errs.Translate(val.Trans)))
		return
	}

	// 业务处理
	resCode, msg := service.AddEnv(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}

// BatchOperation 批量操作[启用/禁用]
func (c *EnvController) BatchOperation(ctx *gin.Context) {
	// 获取参数
	p := new(model.BatchOperationEnv)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		// 判断err是不是validator.ValidationErrors类型
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			res.ResError(ctx, res.CodeInvalidParam)
			return
		}

		// 翻译错误
		res.ResErrorWithMsg(ctx, res.CodeInvalidParam, val.RemoveTopStruct(errs.Translate(val.Trans)))
		return
	}

	// 业务处理
	resCode, msg := service.BatchOperationEnv(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}

// Update 修改
func (c *EnvController) Update(ctx *gin.Context) {
	// 获取参数
	p := new(model.UpdateEnv)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		// 判断err是不是validator.ValidationErrors类型
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			res.ResError(ctx, res.CodeInvalidParam)
			return
		}

		// 翻译错误
		res.ResErrorWithMsg(ctx, res.CodeInvalidParam, val.RemoveTopStruct(errs.Translate(val.Trans)))
		return
	}

	// 业务处理
	resCode, msg := service.UpdateEnv(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}

// BindPanel 绑定面板
func (c *EnvController) BindPanel(ctx *gin.Context) {
	// 获取参数
	p := new(model.BindPanel)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		// 判断err是不是validator.ValidationErrors类型
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			res.ResError(ctx, res.CodeInvalidParam)
			return
		}

		// 翻译错误
		res.ResErrorWithMsg(ctx, res.CodeInvalidParam, val.RemoveTopStruct(errs.Translate(val.Trans)))
		return
	}

	// 业务处理
	resCode, msg := service.BindPanel(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}

// Delete 删除
func (c *EnvController) Delete(ctx *gin.Context) {
	// 获取参数
	p := new(model.DeleteEnv)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		// 判断err是不是validator.ValidationErrors类型
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			res.ResError(ctx, res.CodeInvalidParam)
			return
		}

		// 翻译错误
		res.ResErrorWithMsg(ctx, res.CodeInvalidParam, val.RemoveTopStruct(errs.Translate(val.Trans)))
		return
	}

	// 业务处理
	resCode, msg := service.DeleteEnv(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}
