package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
)

var AllowOrigins = []string{
	"http://127.0.0.1:3000",
	"http://localhost:3000",
}

var CorsConfig = cors.Config{
	AllowAllOrigins:  false,
	AllowOrigins:     AllowOrigins,                                                          // 允许的源，生产环境中应替换为具体的允许域名
	AllowOriginFunc:  func(origin string) bool { return true },                              // 自定义函数来判断源是否允许
	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},  // 允许的HTTP方法列表
	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, // 允许的HTTP头部列表
	AllowCredentials: true,                                                                  // 是否允许浏览器发送Cookie
	MaxAge:           30 * time.Minute,                                                      // 预检请求（OPTIONS）的缓存时间（秒）
}
