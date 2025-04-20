package post

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"example/internal/pkg/core"
	"example/internal/pkg/errno"
	"example/internal/pkg/known"
	"example/internal/pkg/log"
	"example/internal/pkg/model"
)

// CreateRequest 创建文章请求结构
type CreateRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=100"`
	Content string `json:"content" binding:"required,min=1"`
}

// Create 处理创建博客文章请求
func (ctrl *Controller) Create(c *gin.Context) {
	log.Infow("Create post function called")

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	// 获取当前用户信息
	username := c.GetString(known.XUsernameKey)
	if username == "" {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	// 生成唯一的文章ID
	postID := uuid.New().String()

	// 创建文章对象
	post := &model.PostM{
		Username:  username,
		PostID:    postID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	if err := ctrl.postRepo.Create(c, post); err != nil {
		log.Errorw("Failed to create post", "err", err)
		core.WriteResponse(c, errno.ErrDatabase, nil)
		return
	}

	// 为当前用户添加对这篇文章的所有权限
	resourcePath := "/v1/posts/" + postID
	if _, err := ctrl.authz.AddPolicy(username, resourcePath, "*"); err != nil {
		log.Warnw("Failed to add policy for post", "username", username, "postID", postID, "err", err)
	}

	core.WriteResponse(c, nil, gin.H{
		"post_id": postID,
	})
}
