package user

import (
	"example/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

// User 表示用户数据结构
type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// Create 创建新用户
func Create(c *gin.Context) {
	var user User

	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&user); err != nil {
		// 使用错误处理系统返回错误信息
		errno.JsonError(c, errno.ErrBind.WithErr(err))
		return
	}

	// 验证用户信息
	if len(user.Username) == 0 {
		errno.JsonError(c, errno.ErrValidation.WithMsg("用户名不能为空"))
		return
	}

	if len(user.Password) < 6 {
		errno.JsonError(c, errno.ErrValidation.WithMsg("密码长度不能少于6个字符"))
		return
	}

	// 模拟成功创建用户
	errno.WriteResponse(c, nil, map[string]string{
		"username": user.Username,
		"message":  "用户创建成功",
	})
}

// Get 获取用户信息
func Get(c *gin.Context) {
	username := c.Param("username")

	// 检查用户名
	if username == "" {
		errno.JsonError(c, errno.ErrValidation.WithMsg("用户名不能为空"))
		return
	}

	// 模拟用户不存在的情况
	if username == "unknown" {
		errno.JsonError(c, errno.ErrUserNotFound)
		return
	}

	// 模拟成功获取用户
	user := User{
		Username: username,
		Email:    username + "@example.com",
	}

	errno.WriteResponse(c, nil, user)
}

// Delete 删除用户
func Delete(c *gin.Context) {
	username := c.Param("username")

	// 检查用户名
	if username == "" {
		errno.JsonError(c, errno.ErrValidation.WithMsg("用户名不能为空"))
		return
	}

	// 模拟需要管理员权限
	if username == "admin" {
		errno.JsonError(c, errno.ErrForbidden.WithMsg("不能删除管理员用户"))
		return
	}

	// 模拟成功删除用户
	errno.WriteResponse(c, nil, "用户已删除")
}

// Update 更新用户信息
func Update(c *gin.Context) {
	var user User
	username := c.Param("username")

	// 检查用户名
	if username == "" {
		errno.JsonError(c, errno.ErrValidation.WithMsg("用户名不能为空"))
		return
	}

	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&user); err != nil {
		errno.JsonError(c, errno.ErrBind.WithErr(err))
		return
	}

	// 模拟服务器内部错误
	if username == "error" {
		errno.JsonError(c, errno.ErrDatabase.WithMsg("数据库连接失败"))
		return
	}

	// 模拟成功更新用户
	errno.WriteResponse(c, nil, map[string]string{
		"username": username,
		"message":  "用户信息已更新",
	})
}
