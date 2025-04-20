package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"example/internal/pkg/model"
)

// PostRepository 定义文章仓储接口
type PostRepository interface {
	Create(ctx context.Context, post *model.PostM) error
	Get(ctx context.Context, postID string) (*model.PostM, error)
	Update(ctx context.Context, post *model.PostM) error
	Delete(ctx context.Context, postID string) error
	ListByUser(ctx context.Context, username string, offset, limit int) ([]*model.PostM, int64, error)
	List(ctx context.Context, offset, limit int) ([]*model.PostM, int64, error)
}

// postRepo 实现文章仓储接口
type postRepo struct {
	db *gorm.DB
}

// NewPostRepository 创建文章仓储实例
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepo{
		db: db,
	}
}

// Create 创建文章
func (r *postRepo) Create(ctx context.Context, post *model.PostM) error {
	return r.db.WithContext(ctx).Create(post).Error
}

// Get 根据文章ID获取文章
func (r *postRepo) Get(ctx context.Context, postID string) (*model.PostM, error) {
	var post model.PostM
	if err := r.db.WithContext(ctx).Where("postID = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &post, nil
}

// Update 更新文章
func (r *postRepo) Update(ctx context.Context, post *model.PostM) error {
	return r.db.WithContext(ctx).Save(post).Error
}

// Delete 删除文章
func (r *postRepo) Delete(ctx context.Context, postID string) error {
	return r.db.WithContext(ctx).Where("postID = ?", postID).Delete(&model.PostM{}).Error
}

// ListByUser 列出指定用户的文章
func (r *postRepo) ListByUser(ctx context.Context, username string, offset, limit int) ([]*model.PostM, int64, error) {
	var posts []*model.PostM
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.PostM{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("username = ?", username).Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

// List 列出所有文章
func (r *postRepo) List(ctx context.Context, offset, limit int) ([]*model.PostM, int64, error) {
	var posts []*model.PostM
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.PostM{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}
