package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"example/internal/pkg/known"
)

func TestAuthzByRole(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 测试不同角色访问不同权限级别的路由
	testCases := []struct {
		name          string
		userRole      string
		requiredRole  string
		expectedCode  int
		expectedError bool
	}{
		{
			name:          "管理员访问管理员路由",
			userRole:      "admin",
			requiredRole:  "admin",
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "编辑访问编辑路由",
			userRole:      "editor",
			requiredRole:  "editor",
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "读者访问读者路由",
			userRole:      "reader",
			requiredRole:  "reader",
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "管理员访问编辑路由",
			userRole:      "admin",
			requiredRole:  "editor",
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "管理员访问读者路由",
			userRole:      "admin",
			requiredRole:  "reader",
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "编辑访问读者路由",
			userRole:      "editor",
			requiredRole:  "reader",
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name:          "读者访问编辑路由",
			userRole:      "reader",
			requiredRole:  "editor",
			expectedCode:  http.StatusForbidden,
			expectedError: true,
		},
		{
			name:          "编辑访问管理员路由",
			userRole:      "editor",
			requiredRole:  "admin",
			expectedCode:  http.StatusForbidden,
			expectedError: true,
		},
		{
			name:          "读者访问管理员路由",
			userRole:      "reader",
			requiredRole:  "admin",
			expectedCode:  http.StatusForbidden,
			expectedError: true,
		},
		{
			name:          "无角色访问读者路由",
			userRole:      "",
			requiredRole:  "reader",
			expectedCode:  http.StatusForbidden,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个测试路由器
			r := gin.New()

			// 使用角色授权中间件
			r.GET("/test", AuthzByRole(tc.requiredRole), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "authorized"})
			})

			// 创建请求
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// 设置上下文中的用户角色
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			ctx.Set(known.XRoleKey, tc.userRole)

			// 手动调用中间件和处理器
			r.HandleContext(ctx)

			// 验证结果
			assert.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedError {
				assert.Contains(t, w.Body.String(), "code")
				assert.Contains(t, w.Body.String(), "message")
			} else {
				assert.Contains(t, w.Body.String(), "authorized")
			}
		})
	}
}

func TestOnlyOwner(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个测试路由器
	r := gin.New()

	// 使用所有者验证中间件
	r.GET("/users/:username", OnlyOwner(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authorized"})
	})

	// 测试用户访问自己的资源
	req := httptest.NewRequest("GET", "/users/alice", nil)
	w := httptest.NewRecorder()

	// 设置上下文中的用户名
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Set(known.XUsernameKey, "alice")
	ctx.Params = gin.Params{{Key: "username", Value: "alice"}}

	// 手动调用中间件和处理器
	r.HandleContext(ctx)

	// 验证结果 - 应该成功
	assert.Equal(t, http.StatusOK, w.Code)

	// 测试用户访问他人的资源
	req = httptest.NewRequest("GET", "/users/bob", nil)
	w = httptest.NewRecorder()

	// 设置上下文中的用户名
	ctx, _ = gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Set(known.XUsernameKey, "alice")
	ctx.Params = gin.Params{{Key: "username", Value: "bob"}}

	// 手动调用中间件和处理器
	r.HandleContext(ctx)

	// 验证结果 - 应该失败
	assert.Equal(t, http.StatusForbidden, w.Code)
}
