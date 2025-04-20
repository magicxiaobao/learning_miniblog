package schema

import (
	"example/internal/pkg/log"
	"example/internal/pkg/model"

	"gorm.io/gorm"
)

// InitTables 初始化数据库表，自动迁移模型到数据库
func InitTables(db *gorm.DB) error {
	log.Infow("Starting database schema initialization...")

	models := []interface{}{
		&model.UserM{},
		&model.PostM{},
	}

	// 挨个迁移模型
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Errorw("Failed to migrate model", "model", model, "err", err)
			return err
		}
		log.Infow("Model migrated successfully", "model", model)
	}

	log.Infow("All database tables initialized successfully")
	return nil
}

// CreateInitialData 创建初始数据，如默认管理员账户等
func CreateInitialData(db *gorm.DB) error {
	log.Infow("Checking initial data...")

	// 检查是否已存在管理员账户
	var count int64
	if err := db.Model(&model.UserM{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		log.Errorw("Failed to check admin users", "err", err)
		return err
	}

	// 已存在管理员账户，无需创建
	if count > 0 {
		log.Infow("Admin users already exist, skipping creation")
		return nil
	}

	log.Infow("Initial data creation completed")
	return nil
}
