package post

import (
	"example/internal/pkg/auth"
	"example/internal/pkg/repository"

	"gorm.io/gorm"
)

// Controller 博客文章控制器
type Controller struct {
	db       *gorm.DB
	authz    *auth.Authz
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

// New 创建博客文章控制器
func New(db *gorm.DB, authz *auth.Authz) *Controller {
	return &Controller{
		db:       db,
		authz:    authz,
		postRepo: repository.NewPostRepository(db),
		userRepo: repository.NewUserRepository(db),
	}
}
