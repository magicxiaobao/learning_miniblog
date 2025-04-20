package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example/internal/pkg/auth"
	"example/internal/pkg/core"
	"example/internal/pkg/errno"
	"example/internal/pkg/known"
	"example/internal/pkg/log"
)

// Authz 授权中间件，基于Casbin实现访问控制
func Authz(a *auth.Authz) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取资源、操作方法和用户
		obj := c.Request.URL.Path
		act := c.Request.Method
		sub := c.GetString(known.XUsernameKey)

		// 记录授权请求
		log.Infow("Checking authorization", "subject", sub, "object", obj, "action", act)

		// 检查用户是否有权限
		allowed, err := a.Authorize(sub, obj, act)
		if err != nil {
			core.WriteResponse(c, errno.ErrInternalServer.WithErr(err), nil)
			c.Abort()
			return
		}

		if !allowed {
			core.WriteResponse(c, errno.ErrForbidden, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthzByRole 基于角色的授权中间件，检查用户是否拥有指定角色
func AuthzByRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		role := c.GetString(known.XRoleKey)
		if role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    errno.ErrForbidden.Code,
				"message": "用户未登录或无法获取角色信息",
			})
			return
		}

		// 检查用户角色是否满足要求
		if !hasPermission(role, requiredRole) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    errno.ErrForbidden.Code,
				"message": "权限不足，需要" + requiredRole + "角色",
			})
			return
		}

		c.Next()
	}
}

// hasPermission 检查用户角色是否有足够权限
func hasPermission(userRole, requiredRole string) bool {
	// 角色层级: admin > editor > reader
	roleWeight := map[string]int{
		"admin":  100,
		"editor": 50,
		"reader": 10,
		"":       0,
	}

	// 检查用户角色权重是否大于等于所需权重
	return roleWeight[userRole] >= roleWeight[requiredRole]
}

// OnlyOwner 限制只有资源所有者才能访问
func OnlyOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前登录用户
		username := c.GetString(known.XUsernameKey)

		// 获取资源所有者（通常从URL参数或请求体中获取）
		resourceOwner := c.Param("username")

		// 如果资源所有者为空，则尝试从其他地方获取
		if resourceOwner == "" {
			// 这里可以尝试从请求体、查询参数等处获取
			// 例如：resourceOwner = c.Query("owner")
		}

		// 检查当前用户是否为资源所有者
		if username != resourceOwner {
			core.WriteResponse(c, errno.ErrForbidden.WithMsg("只有资源所有者才能执行此操作"), nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
