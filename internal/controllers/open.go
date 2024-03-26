package controllers

import (
	"github.com/gin-gonic/gin"

	"QLToolsV2/internal/model"
	"QLToolsV2/internal/service"
	res "QLToolsV2/pkg/response"
)

type OpenController struct{}

// Router 注册路由
func (c *OpenController) Router(r *gin.RouterGroup) {
	// 登录
	r.POST("login", c.Login)
	// 注册
	r.POST("register", c.Register)
}

// Login 用户登录
func (c *OpenController) Login(ctx *gin.Context) {
	// 获取参数
	p := new(model.Login)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		res.ResErrorWithMsg(ctx, res.CodeInvalidParam, err)
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
		res.ResErrorWithMsg(ctx, res.CodeInvalidParam, err)
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
