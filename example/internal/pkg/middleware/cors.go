package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Cors 实现跨域资源共享 (CORS) 中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置 CORS 相关的响应头
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 计算请求处理时间
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 记录响应时间
		c.Writer.Header().Set("X-Response-Time", latency.String())

		// 继续处理请求
		c.Next()
	}
}
