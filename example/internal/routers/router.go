package routers

import (
	"github.com/gin-gonic/gin"

	"example/internal/routers/api/log"
	"example/internal/routers/api/user"
)

// InstallRouters 安装应用的路由
func InstallRouters(g *gin.Engine) error {
	// 注册 404 Handler
	g.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "Page not found"})
	})

	// 注册 /healthz 路由
	g.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 创建 v1 版本的路由组
	v1 := g.Group("/v1")
	{
		// 日志演示路由
		logGroup := v1.Group("/logs")
		{
			logGroup.GET("/debug", log.DemoDebug)
			logGroup.GET("/info", log.DemoInfo)
			logGroup.GET("/warn", log.DemoWarn)
			logGroup.GET("/error", log.DemoError)
			logGroup.GET("/error-stack", log.DemoErrorWithStack)
		}

		// 用户相关路由
		userGroup := v1.Group("/users")
		{
			// 创建用户
			userGroup.POST("", user.Create)
			// 获取用户信息
			userGroup.GET("/:username", user.Get)
			// 更新用户
			userGroup.PUT("/:username", user.Update)
			// 删除用户
			userGroup.DELETE("/:username", user.Delete)
		}
	}

	return nil
}
