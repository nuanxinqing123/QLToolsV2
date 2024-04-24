package controllers

import (
	"github.com/gin-gonic/gin"

	res "QLToolsV2/pkg/response"
)

type UserController struct{}

// Router 注册路由
func (c *UserController) Router(r *gin.RouterGroup) {
	// 登出
	r.POST("/logout", c.Logout)

	// 获取用户信息
	r.GET("/info", c.Info)
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
