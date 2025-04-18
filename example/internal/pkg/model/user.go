package model

import "time"

// UserM 定义用户模型
type UserM struct {
	ID        uint64    `gorm:"column:id;primary_key"` // 用户ID
	Username  string    `gorm:"column:username;not null"`
	Password  string    `gorm:"column:password;not null"`
	Nickname  string    `gorm:"column:nickname"`
	Email     string    `gorm:"column:email"`
	Phone     string    `gorm:"column:phone"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

// TableName 用户表名
func (u *UserM) TableName() string {
	return "user"
}
