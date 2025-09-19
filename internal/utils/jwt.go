package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bluele/gcache"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	_const "github.com/nuanxinqing123/QLToolsV2/internal/const"
)

// JWT配置常量
const (
	TokenExpiration  = 24 * time.Hour     // Token过期时间：24小时
	RefreshTokenExp  = 7 * 24 * time.Hour // 刷新Token过期时间：7天
	TokenCachePrefix = "jwt:token:"       // Cache中Token的前缀
)

// JWTClaims JWT载荷结构
type JWTClaims struct {
	UserID    int64  `json:"user_id"`    // 用户ID
	TokenType string `json:"token_type"` // Token类型：access/refresh
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	cache gcache.Cache
}

// NewJWTManager 创建JWT管理器实例
func NewJWTManager() *JWTManager {
	return &JWTManager{
		cache: config.Cache,
	}
}

// GenerateTokenPair 生成访问Token和刷新Token对
func (j *JWTManager) GenerateTokenPair(userID int64) (accessToken, refreshToken string, err error) {
	now := time.Now()

	// 生成访问Token
	accessClaims := &JWTClaims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    _const.JWTIssuer,
			Subject:   strconv.FormatInt(userID, 10),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(config.Config.App.Secret))
	if err != nil {
		return "", "", fmt.Errorf("生成访问token失败: %w", err)
	}

	// 生成刷新Token
	refreshClaims := &JWTClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenExp)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    _const.JWTIssuer,
			Subject:   strconv.FormatInt(userID, 10),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(config.Config.App.Secret))
	if err != nil {
		return "", "", fmt.Errorf("生成刷新token失败: %w", err)
	}

	// 将Token存储到Cache中，用于主动注销管理
	accessKey := TokenCachePrefix + "access:" + strconv.FormatInt(userID, 10)
	refreshKey := TokenCachePrefix + "refresh:" + strconv.FormatInt(userID, 10)

	// 存储访问Token（设置过期时间）
	if err = j.cache.SetWithExpire(accessKey, accessToken, TokenExpiration); err != nil {
		return "", "", fmt.Errorf("存储访问token到Cache失败: %w", err)
	}

	// 存储刷新Token（设置过期时间）
	if err = j.cache.SetWithExpire(refreshKey, refreshToken, RefreshTokenExp); err != nil {
		return "", "", fmt.Errorf("存储刷新token到Cache失败: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析Token并验证
func (j *JWTManager) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(config.Config.App.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析token失败: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("无效的token")
	}

	// 检查Token是否在Cache中存在（用于验证是否已被注销）
	var cacheKey string
	if claims.TokenType == "access" {
		cacheKey = TokenCachePrefix + "access:" + strconv.FormatInt(claims.UserID, 10)
	} else {
		cacheKey = TokenCachePrefix + "refresh:" + strconv.FormatInt(claims.UserID, 10)
	}

	storedToken, err := j.cache.Get(cacheKey)
	if err != nil {
		return nil, fmt.Errorf("验证token状态失败: %w", err)
	}

	// 验证Token是否匹配
	if storedToken != tokenString {
		return nil, fmt.Errorf("token不匹配，可能已被替换")
	}

	return claims, nil
}

// RevokeToken 注销Token（删除Cache中的记录）
func (j *JWTManager) RevokeToken(userID int64) error {
	// 删除访问Token和刷新Token
	accessKey := TokenCachePrefix + "access:" + strconv.FormatInt(userID, 10)
	refreshKey := TokenCachePrefix + "refresh:" + strconv.FormatInt(userID, 10)

	j.cache.Remove(accessKey)
	j.cache.Remove(refreshKey)

	return nil
}

// RefreshAccessToken 使用刷新Token生成新的访问Token
func (j *JWTManager) RefreshAccessToken(refreshToken string) (newAccessToken string, err error) {
	// 解析刷新Token
	claims, err := j.ParseToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("刷新token无效: %w", err)
	}

	// 验证是否为刷新Token
	if claims.TokenType != "refresh" {
		return "", fmt.Errorf("提供的不是刷新token")
	}

	// 生成新的访问Token
	now := time.Now()
	newClaims := &JWTClaims{
		UserID:    claims.UserID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    _const.JWTIssuer,
			Subject:   strconv.FormatInt(claims.UserID, 10),
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newAccessToken, err = tokenObj.SignedString([]byte(config.Config.App.Secret))
	if err != nil {
		return "", fmt.Errorf("生成新访问token失败: %w", err)
	}

	// 更新Cache中的访问Token
	accessKey := TokenCachePrefix + "access:" + strconv.FormatInt(claims.UserID, 10)
	if err = j.cache.SetWithExpire(accessKey, newAccessToken, TokenExpiration); err != nil {
		return "", fmt.Errorf("更新Cache中的访问token失败: %w", err)
	}

	return newAccessToken, nil
}
