package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/utils"
)

const (
	// ContextClaimsKey 用于在上下文中存储JWTClaims的键名，便于后续业务处理
	ContextClaimsKey = "jwt_claims"
	// ContextUserIDKey 用于在上下文中缓存当前登录用户ID，减少重复解析
	ContextUserIDKey = "jwt_user_id"
)

// JWTAuth 校验请求头中的JWT访问Token，并根据缓存判断是否已注销
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取Authorization头，标准格式应为Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "缺少Authorization头部信息",
			})
			return
		}

		// 按空格拆分，Bearer前缀与Token内容需要同时存在
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization头格式错误，应为Bearer token",
			})
			return
		}

		// 去除首尾空格，避免因意外空白字符导致验证失败
		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Token内容为空",
			})
			return
		}

		// 使用统一的JWT管理器完成解析与缓存校验
		manager := utils.NewJWTManager()
		claims, err := manager.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Token验证失败，可能已过期或被注销",
			})
			return
		}

		// 限制只允许访问Token通过，刷新Token不可直接访问业务接口
		if claims.TokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "仅支持访问Token访问受保护资源",
			})
			return
		}

		// 将关键身份信息写入上下文，后续处理链可直接复用
		c.Set(ContextClaimsKey, claims)
		c.Set(ContextUserIDKey, claims.UserID)

		// Token校验通过，继续执行后续处理逻辑
		c.Next()
	}
}
