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

type UserController struct{}

// Router 注册路由
func (c *UserController) Router(r *gin.RouterGroup) {
	// 登出
	r.POST("/logout", c.Logout)

	// 获取用户信息
	r.GET("/info", c.Info)
	// 获取登录信息
	r.GET("/login/info", c.LoginInfo)
}

// Logout 登出
func (c *UserController) Logout(ctx *gin.Context) {
	// 删除Cookie
	res.ResSuccess(ctx, "退出成功")
}

// Info 获取用户信息
func (c *UserController) Info(ctx *gin.Context) {
	res.ResSuccess(ctx, "success") // 成功
}

// LoginInfo 获取登录信息
func (c *UserController) LoginInfo(ctx *gin.Context) {
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
	resCode, msg := service.LoginInfo(p)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}
