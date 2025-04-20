package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"example/internal/pkg/known"
	"example/internal/pkg/log"
)

// DemoHandler 是一个演示中间件，展示中间件的工作流程
func DemoHandler() gin.HandlerFunc {
	// 初始化阶段 - 只在应用启动时执行一次
	startTime := time.Now()
	log.Infow("Demo中间件已初始化", "startTime", startTime)

	// 以下是请求处理阶段 - 对每个请求执行
	return func(c *gin.Context) {
		// 1. 请求前的处理
		// 1.1 记录请求开始时间
		requestStartTime := time.Now()

		// 1.2 记录进入中间件的日志
		log.Infow("请求进入Demo中间件",
			"request_id", c.GetString(known.XRequestIDKey),
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
		)

		// 1.3 添加自定义头信息
		c.Header("X-Demo-Middleware", "true")

		// 1.4 向上下文添加信息，供后续中间件或处理器使用
		c.Set("demoKey", "demoValue")

		// 1.5 检查是否要终止处理链
		if c.Query("abort") == "true" {
			// 如果URL中包含 ?abort=true，则终止处理
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "请求被Demo中间件拦截",
			})
			return
		}

		// 2. 调用下一个中间件或处理器
		// 这一行是中间件的关键部分，它将控制权传递给链中的下一个函数
		c.Next()

		// 3. 请求后的处理 (路由处理完成后执行)
		// 3.1 计算请求处理时间
		latency := time.Since(requestStartTime)

		// 3.2 记录离开中间件的日志
		log.Infow("请求离开Demo中间件",
			"request_id", c.GetString(known.XRequestIDKey),
			"status", c.Writer.Status(),
			"latency", latency,
		)

		// 3.3 可以修改响应
		// 注意：此时响应体可能已经发送，不是所有修改都有效
		c.Header("X-Demo-Response-Time", latency.String())
	}
}

// ShowMiddlewareDemo 展示一个特别的路由，用于演示中间件流程
func ShowMiddlewareDemo(r *gin.Engine) {
	// 创建演示用的路由组
	demo := r.Group("/demo")

	// 在路由组上应用中间件
	demo.Use(DemoHandler())

	// 定义路由处理函数
	demo.GET("/middleware", func(c *gin.Context) {
		// 从上下文获取中间件设置的值
		demoValue, exists := c.Get("demoKey")

		// 返回响应
		c.JSON(http.StatusOK, gin.H{
			"message":    "这是一个演示中间件流程的端点",
			"demoExists": exists,
			"demoValue":  demoValue,
			"headers":    c.Request.Header,
		})
	})

	// 在单个路由上使用额外的中间件
	demo.GET("/middleware/stacked", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "这个路由使用了多层中间件",
		})
	}, func(c *gin.Context) {
		// 这是直接在路由定义中添加的中间件
		log.Infow("路由内联中间件被调用",
			"path", c.Request.URL.Path,
		)
	})

	// 演示中间件终止流的路由
	demo.GET("/middleware/abort", func(c *gin.Context) {
		// 这个处理器永远不会被调用，因为中间件会拦截包含 ?abort=true 的请求
		c.JSON(http.StatusOK, gin.H{
			"message": "如果你看到这个消息，意味着中间件没有拦截请求",
		})
	})
}
