package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"example/internal/pkg/model"
)

// UserRepository 定义用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *model.UserM) error
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	Delete(ctx context.Context, username string) error
	List(ctx context.Context, offset, limit int) ([]*model.UserM, int64, error)
	GetByID(ctx context.Context, id uint64) (*model.UserM, error)
}

// userRepo 实现用户仓储接口
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

// Create 创建用户
func (r *userRepo) Create(ctx context.Context, user *model.UserM) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Get 根据用户名获取用户
func (r *userRepo) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// GetByID 根据ID获取用户
func (r *userRepo) GetByID(ctx context.Context, id uint64) (*model.UserM, error) {
	var user model.UserM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepo) Update(ctx context.Context, user *model.UserM) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r *userRepo) Delete(ctx context.Context, username string) error {
	return r.db.WithContext(ctx).Where("username = ?", username).Delete(&model.UserM{}).Error
}

// List 列出用户
func (r *userRepo) List(ctx context.Context, offset, limit int) ([]*model.UserM, int64, error) {
	var users []*model.UserM
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.UserM{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}
