package middleware

import (
	"errors"
	"time"

	"github.com/bluele/gcache"
	"github.com/gin-gonic/gin"

	"QLToolsV2/config"
	res "QLToolsV2/pkg/response"
)

// IdempotentConfig 幂等配置
type IdempotentConfig struct {
	// Token在请求头中的键名
	TokenHeader string
	// 缓存过期时间
	ExpireTime time.Duration
	// 是否在处理完成后自动删除Token
	AutoDelete bool
}

// DefaultIdempotentConfig 默认幂等配置
var DefaultIdempotentConfig = IdempotentConfig{
	TokenHeader: "Idempotent-Token",
	ExpireTime:  time.Hour,
	AutoDelete:  true,
}

// Idempotent 幂等中间件
func Idempotent(cfg ...IdempotentConfig) gin.HandlerFunc {
	// 使用默认配置或自定义配置
	defaultConfig := DefaultIdempotentConfig
	if len(cfg) > 0 {
		defaultConfig = cfg[0]
	}

	return func(c *gin.Context) {
		// 跳过GET请求
		if c.Request.Method != "GET" {
			// 获取幂等Token
			token := c.Request.Header.Get(defaultConfig.TokenHeader)
			if token == "" {
				res.ResErrorWithMsg(c, res.CodeInvalidParam, "缺少Token")
				c.Abort()
				return
			}

			// 生成缓存key
			cacheKey := "idempotent:" + token

			// 检查是否重复请求
			_, err := config.GinCache.Get(cacheKey)
			if !errors.Is(err, gcache.KeyNotFoundError) {
				// 重复请求
				res.ResErrorWithMsg(c, res.CodeInvalidParam, "请勿重复请求")
				c.Abort()
				return
			}

			// 写入缓存
			err = config.GinCache.SetWithExpire(cacheKey, true, defaultConfig.ExpireTime)
			if err != nil {
				config.GinLOG.Error(err.Error())
				res.ResError(c, res.CodeServerBusy)
				c.Abort()
				return
			}

			// 处理请求
			c.Next()

			// 如果配置了自动删除，则在请求处理完成后删除Token
			if defaultConfig.AutoDelete {
				_ = config.GinCache.Remove(cacheKey)
			}
		}
	}
}
