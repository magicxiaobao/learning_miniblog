package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example/internal/api/post"
	"example/internal/api/user"
	"example/internal/pkg/auth"
	"example/internal/pkg/middleware"
)

// Router 包含所有路由定义
type Router struct {
	db       *gorm.DB
	authz    *auth.Authz
	userCtrl *user.Controller
	postCtrl *post.Controller
}

// New 创建路由器实例
func New(db *gorm.DB, authz *auth.Authz) *Router {
	return &Router{
		db:       db,
		authz:    authz,
		userCtrl: user.New(db, authz),
		postCtrl: post.New(db, authz),
	}
}

// Load 加载所有路由
func (r *Router) Load(g *gin.Engine) {
	// 中间件
	g.Use(middleware.RequestID())
	g.Use(middleware.Logger())
	g.Use(middleware.Recovery())
	g.Use(middleware.Cors())
	g.Use(middleware.Secure())
	g.Use(middleware.NoCache())

	// 健康检查
	g.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API 版本 v1
	v1 := g.Group("/v1")
	{
		// 无需认证的接口
		// 用户注册
		v1.POST("/users", r.userCtrl.Create)
		// 用户登录
		v1.POST("/login", r.userCtrl.Login)

		// 需要认证的接口
		auth := v1.Group("")
		auth.Use(middleware.Authn())
		{
			// 用户登出
			auth.POST("/logout", r.userCtrl.Logout)

			// 用户相关路由
			users := auth.Group("/users")
			{
				// 获取特定用户信息 - 使用 OnlyOwner 中间件确保只有本人或管理员可以操作
				users.GET("/:username", r.userCtrl.Get)

				// 以下是需要通过 Casbin 授权的路由示例
				// 用 authz 中间件保护的接口
				admin := users.Group("")
				admin.Use(middleware.AuthzByRole("admin"))
				{
					// 这里可以添加仅管理员可访问的用户相关接口
					// 例如：获取所有用户列表、删除用户等
				}

				// 用户相关的文章路由
				users.GET("/:username/posts", r.postCtrl.ListByUser)
			}

			// 博客文章相关路由
			posts := auth.Group("/posts")
			{
				// 创建文章
				posts.POST("", r.postCtrl.Create)

				// 获取指定文章
				posts.GET("/:postID", r.postCtrl.Get)

				// 获取文章列表
				posts.GET("", r.postCtrl.List)

				// 添加更多文章相关API
				// posts.PUT("/:postID", r.postCtrl.Update)
				// posts.DELETE("/:postID", r.postCtrl.Delete)
			}
		}
	}
}

// API 规范:
// GET    /v1/users       - 获取用户列表
// POST   /v1/users       - 创建用户
// GET    /v1/users/:name - 获取指定用户详情
// PUT    /v1/users/:name - 更新指定用户
// DELETE /v1/users/:name - 删除指定用户

// GET    /v1/posts       - 获取文章列表
// POST   /v1/posts       - 创建文章
// GET    /v1/posts/:id   - 获取指定文章详情
// PUT    /v1/posts/:id   - 更新指定文章
// DELETE /v1/posts/:id   - 删除指定文章
