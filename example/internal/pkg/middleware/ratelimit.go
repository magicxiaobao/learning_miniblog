package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"example/internal/pkg/errno"
	"example/internal/pkg/known"
)

// 保存每个IP的限流器
var (
	// ipLimiters 存储每个 IP 的限流器
	ipLimiters = make(map[string]*rate.Limiter)
	// mu 保护 ipLimiters 的并发访问
	mu sync.Mutex
)

// 获取特定IP的限流器
func getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := ipLimiters[ip]
	if !exists {
		limiter = rate.NewLimiter(r, b)
		ipLimiters[ip] = limiter
	}

	return limiter
}

// RateLimiter 返回基于IP的限流中间件
// r: 每秒允许的请求数
// b: 突发请求的最大数量
func RateLimiter(r float64, b int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端 IP
		ip := c.ClientIP()

		// 获取此 IP 的限流器
		limiter := getLimiter(ip, rate.Limit(r), b)

		// 尝试获取令牌
		if !limiter.Allow() {
			// 如果没有令牌，返回 429 错误
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    errno.ErrTooManyRequests.Code,
				"message": "请求频率超过限制，请稍后再试",
			})
			return
		}

		// 在响应头中添加限流信息
		c.Header("X-RateLimit-Limit", "60")
		c.Header("X-RateLimit-Remaining", "59") // 简化处理

		// 将 IP 地址添加到上下文，便于日志记录
		c.Set(known.XRealIPKey, ip)

		c.Next()
	}
}

// PathRateLimiter 返回基于路径的限流中间件
// 不同的路径可以有不同的限流规则
func PathRateLimiter() gin.HandlerFunc {
	// 路径限流配置
	pathLimits := map[string]struct {
		rate  float64
		burst int
	}{
		"/v1/users": {5, 10},  // 用户相关接口限流较严格
		"/v1/posts": {10, 20}, // 博客接口限流较宽松
		"/v1/logs":  {20, 40}, // 日志接口限流更宽松
	}

	// 每个路径+IP组合的限流器
	limiters := make(map[string]*rate.Limiter)
	var pathMu sync.Mutex

	return func(c *gin.Context) {
		// 获取客户端 IP 和请求路径
		ip := c.ClientIP()
		path := c.FullPath()

		// 寻找匹配的路径限流配置
		var r float64 = 1 // 默认限流：每秒1个请求
		var b int = 5     // 默认突发请求数：5

		// 检查是否有特定路径的配置
		for pattern, limit := range pathLimits {
			if path == pattern || (len(path) > len(pattern) && path[:len(pattern)] == pattern) {
				r = limit.rate
				b = limit.burst
				break
			}
		}

		// 创建键：IP+路径
		key := ip + ":" + path

		// 获取或创建限流器
		pathMu.Lock()
		limiter, exists := limiters[key]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(r), b)
			limiters[key] = limiter
		}
		pathMu.Unlock()

		// 尝试获取令牌
		if !limiter.Allow() {
			// 如果没有令牌，返回 429 错误
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    errno.ErrTooManyRequests.Code,
				"message": "请求频率超过限制，请稍后再试",
			})
			return
		}

		c.Next()
	}
}
