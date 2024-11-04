package middleware

import (
	"errors"
	"time"

	"github.com/bluele/gcache"
	"github.com/gin-gonic/gin"

	"QLToolsV2/config"
	res "QLToolsV2/pkg/response"
)

// Idempotent 幂等中间件
func Idempotent() gin.HandlerFunc {
	return func(c *gin.Context) {
		Token := c.Request.Header.Get("Token")
		config.GinLOG.Debug(Token)
		if Token == "" {
			// 旧客户端 或者 非法请求
			res.ResError(c, res.CodeInvalidParam)
			c.Abort()
			return
		}

		// 检查是否重复请求
		_, err := config.GinCache.Get(Token)
		if !errors.Is(err, gcache.KeyNotFoundError) {
			// 重复请求
			res.ResSuccess(c, "Idempotent request")
			c.Abort()
			return
		}

		// 写入缓存
		_ = config.GinCache.SetWithExpire(Token, Token, time.Hour)
		c.Next()
	}
}
