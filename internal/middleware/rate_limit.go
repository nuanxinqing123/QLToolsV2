package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nuanxinqing123/QLToolsV2/internal/pkg/response"
)

// TokenBucket 令牌桶结构
type TokenBucket struct {
	capacity   int64      // 桶容量
	tokens     int64      // 当前令牌数
	refillRate int64      // 每秒补充的令牌数
	lastRefill time.Time  // 上次补充时间
	mutex      sync.Mutex // 互斥锁
}

// NewTokenBucket 创建新的令牌桶
func NewTokenBucket(capacity, refillRate int64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity, // 初始时桶是满的
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// TakeToken 尝试从桶中取出一个令牌
func (tb *TokenBucket) TakeToken() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// 计算需要补充的令牌数
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tokensToAdd := int64(elapsed * float64(tb.refillRate))

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}

	// 尝试取出令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// RateLimiter 限速器结构
type RateLimiter struct {
	buckets    map[string]*TokenBucket // IP地址到令牌桶的映射
	mutex      sync.RWMutex            // 读写锁
	capacity   int64                   // 桶容量
	refillRate int64                   // 补充速率
}

// NewRateLimiter 创建新的限速器
func NewRateLimiter(capacity, refillRate int64) *RateLimiter {
	return &RateLimiter{
		buckets:    make(map[string]*TokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

// GetBucket 获取或创建指定IP的令牌桶
func (rl *RateLimiter) GetBucket(ip string) *TokenBucket {
	rl.mutex.RLock()
	bucket, exists := rl.buckets[ip]
	rl.mutex.RUnlock()

	if exists {
		return bucket
	}

	// 如果不存在，创建新的令牌桶
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// 双重检查，防止并发创建
	if bucket, exists := rl.buckets[ip]; exists {
		return bucket
	}

	bucket = NewTokenBucket(rl.capacity, rl.refillRate)
	rl.buckets[ip] = bucket
	return bucket
}

// CleanupExpiredBuckets 清理过期的令牌桶（可选的后台清理任务）
func (rl *RateLimiter) CleanupExpiredBuckets(maxIdle time.Duration) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for ip, bucket := range rl.buckets {
		bucket.mutex.Lock()
		if now.Sub(bucket.lastRefill) > maxIdle {
			delete(rl.buckets, ip)
		}
		bucket.mutex.Unlock()
	}
}

// 全局限速器实例
var (
	// 公开接口限速器：每秒10个请求，桶容量20
	openAPILimiter = NewRateLimiter(20, 10)

	// 提交接口限速器：每秒2个请求，桶容量5（更严格的限制）
	submitAPILimiter = NewRateLimiter(5, 2)
)

// RateLimitMiddleware 通用限速中间件
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取该IP的令牌桶
		bucket := limiter.GetBucket(clientIP)

		// 尝试获取令牌
		if !bucket.TakeToken() {
			// 没有令牌，返回限速错误
			response.ResErrorWithMsg(c, response.CodeTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		// 有令牌，继续处理请求
		c.Next()
	}
}

// OpenAPIRateLimit 公开接口限速中间件
func OpenAPIRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(openAPILimiter)
}

// SubmitAPIRateLimit 提交接口限速中间件（更严格）
func SubmitAPIRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(submitAPILimiter)
}

// StartCleanupTask 启动清理任务（可选）
func StartCleanupTask() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
		defer ticker.Stop()

		for range ticker.C {
			// 清理10分钟未活动的令牌桶
			openAPILimiter.CleanupExpiredBuckets(10 * time.Minute)
			submitAPILimiter.CleanupExpiredBuckets(10 * time.Minute)
		}
	}()
}
