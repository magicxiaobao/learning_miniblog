package user

import (
	"time"

	"github.com/gin-gonic/gin"

	"example/internal/pkg/auth"
	"example/internal/pkg/core"
	"example/internal/pkg/errno"
	"example/internal/pkg/log"
	"example/internal/pkg/model"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=30"`
	Nickname string `json:"nickname" binding:"required,min=1,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Role     string `json:"role" binding:"omitempty,oneof=admin editor reader"` // 角色可选，默认为reader
}

// Create 处理用户注册
func (ctrl *Controller) Create(c *gin.Context) {
	log.Infow("Create user function called")

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	// 检查用户名是否已存在
	var count int64
	if err := ctrl.db.Model(&model.UserM{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		core.WriteResponse(c, errno.ErrDatabase.WithErr(err), nil)
		return
	}

	if count > 0 {
		core.WriteResponse(c, errno.ErrUserExists, nil)
		return
	}

	// 加密密码
	hashedPassword, err := auth.Encrypt(req.Password)
	if err != nil {
		core.WriteResponse(c, errno.ErrInternalServer.WithErr(err), nil)
		return
	}

	// 设置默认角色
	role := req.Role
	if role == "" {
		role = "reader"
	}

	// 创建用户
	user := model.UserM{
		Username:  req.Username,
		Password:  hashedPassword,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Phone:     req.Phone,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ctrl.db.Create(&user).Error; err != nil {
		core.WriteResponse(c, errno.ErrDatabase.WithErr(err), nil)
		return
	}

	// 为用户添加基本权限
	// 用户可以访问自己的资源
	resourcePath := "/v1/users/" + user.Username
	if _, err := ctrl.authz.AddPolicy(user.Username, resourcePath, "*"); err != nil {
		log.Warnw("Failed to add policy for user", "username", user.Username, "err", err)
	}

	// 添加角色
	if _, err := ctrl.authz.AddRoleForUser(user.Username, role); err != nil {
		log.Warnw("Failed to add role for user", "username", user.Username, "role", role, "err", err)
	}

	core.WriteResponse(c, nil, nil)
}
