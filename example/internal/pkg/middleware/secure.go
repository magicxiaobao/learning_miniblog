package middleware

import (
	"github.com/gin-gonic/gin"
)

// Secure 是一个 Gin 中间件，用于设置一些安全相关的 HTTP 头
func Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止浏览器执行嵌入页面中的恶意脚本
		c.Header("X-XSS-Protection", "1; mode=block")

		// 控制如何处理网站的 frame，防止点击劫持
		c.Header("X-Frame-Options", "DENY")

		// 禁止浏览器根据响应内容推断响应类型
		c.Header("X-Content-Type-Options", "nosniff")

		// 限制哪些域可以加载资源
		c.Header("Content-Security-Policy", "default-src 'self'")

		if c.Request.TLS != nil {
			// 要求浏览器在指定的时间内只使用 HTTPS 访问当前域名
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// 继续处理请求
		c.Next()
	}
}
