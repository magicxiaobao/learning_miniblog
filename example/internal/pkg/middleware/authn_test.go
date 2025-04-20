package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"example/internal/pkg/auth"
	"example/internal/pkg/known"
)

func TestAuthn(t *testing.T) {
	// 初始化JWT配置
	auth.InitToken(auth.TokenConfig{
		SigningKey: "test-key",
		ExpireTime: 1, // 1小时
	})

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个测试路由器
	r := gin.New()
	r.Use(Authn())

	// 添加一个测试处理器
	r.GET("/test", func(c *gin.Context) {
		username := c.GetString(known.XUsernameKey)
		userID := c.GetUint64(known.XUserIDKey)
		role := c.GetString(known.XRoleKey)

		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"user_id":  userID,
			"role":     role,
		})
	})

	// 测试无Token情况
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 应该返回未授权错误
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 生成有效Token
	token, err := auth.GenerateToken(1, "testuser", "admin")
	assert.NoError(t, err)

	// 测试有效Token
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 应该返回成功
	assert.Equal(t, http.StatusOK, w.Code)
	// 检查响应内容
	assert.Contains(t, w.Body.String(), "testuser")
	assert.Contains(t, w.Body.String(), "admin")
}

func TestAuthnViaQuery(t *testing.T) {
	// 初始化JWT配置
	auth.InitToken(auth.TokenConfig{
		SigningKey: "test-key",
		ExpireTime: 1, // 1小时
	})

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个测试路由器
	r := gin.New()
	r.Use(AuthnViaQuery())

	// 添加一个测试处理器
	r.GET("/test", func(c *gin.Context) {
		username := c.GetString(known.XUsernameKey)
		userID := c.GetUint64(known.XUserIDKey)
		role := c.GetString(known.XRoleKey)

		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"user_id":  userID,
			"role":     role,
		})
	})

	// 测试无Token情况
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 应该返回未授权错误
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 生成有效Token
	token, err := auth.GenerateToken(1, "testuser", "admin")
	assert.NoError(t, err)

	// 测试有效Token（通过查询参数）
	req = httptest.NewRequest("GET", "/test?token="+token, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 应该返回成功
	assert.Equal(t, http.StatusOK, w.Code)
	// 检查响应内容
	assert.Contains(t, w.Body.String(), "testuser")
	assert.Contains(t, w.Body.String(), "admin")
}
