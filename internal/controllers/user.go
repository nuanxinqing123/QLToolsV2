package controllers

import (
	"github.com/gin-gonic/gin"

	"QLToolsV2/internal/service"
	res "QLToolsV2/pkg/response"
)

const CtxUserID = "UserID"

type UserController struct{}

// Router 注册路由
func (c *UserController) Router(r *gin.RouterGroup) {
	// 登出
	r.POST("logout", c.Logout)

	// 获取用户信息
	r.GET("info", c.Info)
}

// Logout 登出
func (c *UserController) Logout(ctx *gin.Context) {
	// 删除Cookie
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	res.ResSuccess(ctx, "退出成功")
}

// Info 获取用户信息
func (c *UserController) Info(ctx *gin.Context) {
	// 获取 UserID
	userId := ctx.GetString(CtxUserID)

	// 业务处理
	resCode, msg := service.Info(userId)
	if resCode == res.CodeSuccess {
		res.ResSuccess(ctx, msg) // 成功
	} else {
		res.ResErrorWithMsg(ctx, resCode, msg) // 失败
	}
}
