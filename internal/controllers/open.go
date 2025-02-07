package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

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
	// 获取幂等Token
	r.GET("/idempotent/token", c.GetIdempotentToken)

	// KEY检查
	r.POST("/key_check", c.KeyCheck)
	// 在线服务
	r.GET("/online/service", c.OnlineService)
	// 提交服务
	r.POST("/submit/service", c.SubmitService)
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

// SubmitService 提交服务
func (c *OpenController) SubmitService(ctx *gin.Context) {
	// 获取参数
	p := new(model.Submit)
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
	resCode, msg := service.SubmitService(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}

// GetIdempotentToken 获取幂等Token
func (c *OpenController) GetIdempotentToken(ctx *gin.Context) {
	// 生成UUID作为幂等Token
	token := uuid.NewString()

	// 返回Token
	res.ResSuccess(ctx, gin.H{
		"token": token,
	})
}
