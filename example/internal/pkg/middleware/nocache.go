package middleware

import (
	"github.com/gin-gonic/gin"
)

// NoCache 是一个 Gin 中间件，用于禁用客户端缓存
// 适用于希望始终从服务器获取最新数据的 API
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置禁止缓存的 HTTP 头
		c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Header("Last-Modified", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Header("Pragma", "no-cache")

		// 继续处理请求
		c.Next()
	}
}
