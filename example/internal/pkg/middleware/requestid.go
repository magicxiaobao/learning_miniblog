package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"example/internal/pkg/known"
)

// RequestID 是一个中间件，用来在请求上下文中设置唯一请求 ID。
// 如果请求头中已经有 X-Request-ID，则复用该值，否则生成一个新的 UUID。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头是否已有 X-Request-ID
		requestID := c.Request.Header.Get(known.XRequestIDKey)

		// 如果没有，则生成新的 UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将请求 ID 设置到请求上下文中
		c.Set(known.XRequestIDKey, requestID)

		// 添加到响应头
		c.Writer.Header().Set(known.XRequestIDKey, requestID)

		// 继续处理请求
		c.Next()
	}
}
