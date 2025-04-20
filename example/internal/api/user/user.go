package user

import (
	"example/internal/pkg/auth"

	"gorm.io/gorm"
)

// Controller 用户控制器
type Controller struct {
	db    *gorm.DB
	authz *auth.Authz
}

// New 创建用户控制器
func New(db *gorm.DB, authz *auth.Authz) *Controller {
	return &Controller{
		db:    db,
		authz: authz,
	}
}
