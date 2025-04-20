package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example/internal/pkg/auth"
	"example/internal/pkg/core"
	"example/internal/pkg/errno"
	"example/internal/pkg/known"
)

// Authn 是认证中间件，用于验证用户Token并将用户信息设置到上下文中
func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中解析JWT令牌
		claims, err := auth.ParseRequest(c)
		if err != nil {
			// 根据错误类型返回不同的响应
			switch err {
			case auth.ErrMissingHeader:
				core.WriteResponse(c, errno.ErrTokenInvalid.WithMsg("请求未携带Token"), nil)
			case auth.ErrTokenExpired:
				core.WriteResponse(c, errno.ErrTokenExpired, nil)
			default:
				core.WriteResponse(c, errno.ErrTokenInvalid, nil)
			}
			c.Abort()
			return
		}

		// 将用户信息存储到请求上下文中
		c.Set(known.XUsernameKey, claims.Username)
		c.Set(known.XUserIDKey, claims.UserID)
		c.Set(known.XRoleKey, claims.Role)

		c.Next()
	}
}

// AuthnViaQuery 是认证中间件的变体，从查询参数中获取token（用于WebSocket等特殊场景）
func AuthnViaQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从URL查询参数中获取token
		token := c.Query("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    errno.ErrTokenInvalid.Code,
				"message": "请求未携带Token",
			})
			return
		}

		// 解析Token
		claims, err := auth.ParseToken(token)
		if err != nil {
			code := errno.ErrTokenInvalid
			if err == auth.ErrTokenExpired {
				code = errno.ErrTokenExpired
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    code.Code,
				"message": code.Message,
			})
			return
		}

		// 将用户信息存储到请求上下文中
		c.Set(known.XUsernameKey, claims.Username)
		c.Set(known.XUserIDKey, claims.UserID)
		c.Set(known.XRoleKey, claims.Role)

		c.Next()
	}
}
