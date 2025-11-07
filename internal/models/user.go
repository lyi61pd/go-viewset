package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `gorm:"size:100;not null" json:"name" binding:"required"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email" binding:"required,email"`
	Status    string         `gorm:"size:20;default:inactive" json:"status"`
	Age       int            `gorm:"default:0" json:"age"`
	Phone     string         `gorm:"size:20" json:"phone"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
