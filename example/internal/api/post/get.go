package post

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example/internal/pkg/core"
	"example/internal/pkg/errno"
	"example/internal/pkg/known"
	"example/internal/pkg/log"
)

// PostResponse 文章响应结构
type PostResponse struct {
	PostID    string `json:"post_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Get 获取指定博客文章
func (ctrl *Controller) Get(c *gin.Context) {
	log.Infow("Get post function called")

	// 获取文章ID参数
	postID := c.Param("postID")
	if postID == "" {
		core.WriteResponse(c, errno.ErrParam, nil)
		return
	}

	// 从数据库获取文章
	post, err := ctrl.postRepo.Get(c, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			core.WriteResponse(c, errno.ErrPostNotFound, nil)
			return
		}
		log.Errorw("Failed to get post", "postID", postID, "err", err)
		core.WriteResponse(c, errno.ErrDatabase, nil)
		return
	}

	// 获取当前用户信息，检查权限
	username := c.GetString(known.XUsernameKey)
	role := c.GetString(known.XRoleKey)

	// 只有文章作者和管理员可以查看
	if post.Username != username && role != "admin" {
		// 尝试通过Casbin检查权限
		ok, err := ctrl.authz.Authorize(username, "/v1/posts/"+postID, "GET")
		if err != nil || !ok {
			core.WriteResponse(c, errno.ErrForbidden, nil)
			return
		}
	}

	// 构建响应
	resp := &PostResponse{
		PostID:    post.PostID,
		Title:     post.Title,
		Content:   post.Content,
		Username:  post.Username,
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	core.WriteResponse(c, nil, resp)
}
