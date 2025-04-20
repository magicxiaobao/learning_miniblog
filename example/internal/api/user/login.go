package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example/internal/pkg/auth"
	"example/internal/pkg/core"
	"example/internal/pkg/errno"
	"example/internal/pkg/log"
	"example/internal/pkg/model"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}

// Login 处理用户登录
func (ctrl *Controller) Login(c *gin.Context) {
	log.Infow("Login function called")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	// 查询用户
	var user model.UserM
	if err := ctrl.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	// 验证密码
	if err := auth.Compare(user.Password, req.Password); err != nil {
		core.WriteResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// 生成JWT令牌
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		core.WriteResponse(c, errno.ErrInternalServer.WithErr(err), nil)
		return
	}

	core.WriteResponse(c, nil, &LoginResponse{
		Token: token,
	})
}

// Logout 处理用户登出
func (ctrl *Controller) Logout(c *gin.Context) {
	// JWT是无状态的，客户端只需要删除令牌
	// 这里只是为了提供一个API接口
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登出成功",
	})
}
