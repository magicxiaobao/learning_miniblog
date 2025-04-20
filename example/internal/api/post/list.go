package post

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"example/internal/pkg/core"
	"example/internal/pkg/errno"
	"example/internal/pkg/log"
)

// ListRequest 列表请求参数
type ListRequest struct {
	Offset int `form:"offset" binding:"omitempty,min=0"`
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// ListResponse 文章列表响应
type ListResponse struct {
	TotalCount int64           `json:"total_count"`
	Posts      []*PostResponse `json:"posts"`
}

// List 获取文章列表
func (ctrl *Controller) List(c *gin.Context) {
	log.Infow("List posts function called")

	// 解析请求参数
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	// 设置默认值
	if req.Limit == 0 {
		req.Limit = 10
	}

	// 从数据库获取文章列表
	posts, count, err := ctrl.postRepo.List(c, req.Offset, req.Limit)
	if err != nil {
		log.Errorw("Failed to list posts", "err", err)
		core.WriteResponse(c, errno.ErrDatabase, nil)
		return
	}

	// 构建响应
	postList := make([]*PostResponse, 0, len(posts))
	for _, post := range posts {
		postList = append(postList, &PostResponse{
			PostID:    post.PostID,
			Title:     post.Title,
			Content:   post.Content,
			Username:  post.Username,
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	core.WriteResponse(c, nil, &ListResponse{
		TotalCount: count,
		Posts:      postList,
	})
}

// ListByUser 获取指定用户的文章列表
func (ctrl *Controller) ListByUser(c *gin.Context) {
	log.Infow("ListByUser function called")

	// 获取用户名
	username := c.Param("username")
	if username == "" {
		core.WriteResponse(c, errno.ErrParam, nil)
		return
	}

	// 解析分页参数
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 100 {
		limit = 100
	}

	// 从数据库获取用户的文章列表
	posts, count, err := ctrl.postRepo.ListByUser(c, username, offset, limit)
	if err != nil {
		log.Errorw("Failed to list user posts", "username", username, "err", err)
		core.WriteResponse(c, errno.ErrDatabase, nil)
		return
	}

	// 构建响应
	postList := make([]*PostResponse, 0, len(posts))
	for _, post := range posts {
		postList = append(postList, &PostResponse{
			PostID:    post.PostID,
			Title:     post.Title,
			Content:   post.Content,
			Username:  post.Username,
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	core.WriteResponse(c, nil, &ListResponse{
		TotalCount: count,
		Posts:      postList,
	})
}
