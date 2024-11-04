package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"QLToolsV2/config"
	res "QLToolsV2/pkg/response"
)

// Idempotent 幂等中间件
func Idempotent() gin.HandlerFunc {
	return func(c *gin.Context) {
		uToken := c.Request.Header.Get("uToken")
		if uToken == "" {
			// 旧客户端 或者 非法请求
			res.ResError(c, res.CodeInvalidParam)
			return
		}

		// 检查是否重复请求
		token, err := config.GinCache.Get(uToken)
		if err != nil {
			// 旧客户端 或者 非法请求
			res.ResError(c, res.CodeInvalidParam)
			return
		}

		if token != "" {
			// 重复请求
			res.ResSuccess(c, "Idempotent request")
			return
		}

		// 写入缓存
		_ = config.GinCache.SetWithExpire(uToken, uToken, time.Hour)
		c.Next()
	}
}
