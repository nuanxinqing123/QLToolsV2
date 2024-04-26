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

type OpenController struct{}

// Router 注册路由
func (c *OpenController) Router(r *gin.RouterGroup) {
	// 登录
	r.POST("/login", c.Login)
	// 注册
	r.POST("/register", c.Register)

	// KEY检查
	r.POST("/key_check", c.KeyCheck)
	// 在线服务
	r.GET("/online/service", c.OnlineService)
}

// Login 用户登录
func (c *OpenController) Login(ctx *gin.Context) {
	// 获取参数
	p := new(model.Login)
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
	resCode, msg := service.Login(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}

// Register 用户注册
func (c *OpenController) Register(ctx *gin.Context) {
	// 获取参数
	p := new(model.Register)
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
	resCode, msg := service.Register(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}

// KeyCheck KEY检查
func (c *OpenController) KeyCheck(ctx *gin.Context) {
	// 获取参数
	p := new(model.KeyCheck)
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
	resCode, msg := service.KeyCheck(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}

// OnlineService 在线服务
func (c *OpenController) OnlineService(ctx *gin.Context) {
	// 业务处理
	resCode, msg := service.OnlineService()
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}
