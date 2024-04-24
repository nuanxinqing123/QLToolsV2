package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	_const "QLToolsV2/const"
	res "QLToolsV2/pkg/response"
	"QLToolsV2/utils"
)

// Auth 权限认证
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Token放在Header的Authorization中，并使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// 无Token
			res.ResError(c, res.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			res.ResErrorWithMsg(c, res.CodeNeedLogin, "用户状态已失效")
			c.Abort()
			return
		}

		// 初始化 JWT
		jwt := utils.NewJWT()

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			res.ResErrorWithMsg(c, res.CodeNeedLogin, "用户状态已失效")
			c.Abort()
			return
		}

		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(_const.CtxUserID, mc.UserID)
		c.Next()
	}
}
