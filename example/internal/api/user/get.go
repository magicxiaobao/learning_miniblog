package user

import (
	"github.com/gin-gonic/gin"

	"example/internal/pkg/core"
	"example/internal/pkg/errno"
	"example/internal/pkg/known"
	"example/internal/pkg/log"
	"example/internal/pkg/model"
)

// UserResponse 用户信息响应
type UserResponse struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Get 获取用户详情
func (ctrl *Controller) Get(c *gin.Context) {
	log.Infow("Get user function called")

	// 获取路径参数中的用户名
	username := c.Param("username")
	if username == "" {
		core.WriteResponse(c, errno.ErrParam, nil)
		return
	}

	// 检查当前用户是否有权限获取该用户信息
	// 当前用户只能查看自己的信息，除非是管理员
	currentUser := c.GetString(known.XUsernameKey)
	currentRole := c.GetString(known.XRoleKey)

	if currentUser != username && currentRole != "admin" {
		core.WriteResponse(c, errno.ErrForbidden, nil)
		return
	}

	// 查询用户
	var user model.UserM
	if err := ctrl.db.Where("username = ?", username).First(&user).Error; err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	// 构建响应
	resp := &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	core.WriteResponse(c, nil, resp)
}
