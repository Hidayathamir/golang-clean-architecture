package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	LikeTableName       = "likes"
	LikeColumnID        = "id"
	LikeColumnUserID    = "user_id"
	LikeColumnImageID   = "image_id"
	LikeColumnCreatedAt = "created_at"
	LikeColumnUpdatedAt = "updated_at"
	LikeColumnDeletedAt = "deleted_at"
)

type Like struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	UserID    int64          `gorm:"column:user_id"`
	ImageID   int64          `gorm:"column:image_id"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	User  User  `gorm:"foreignKey:user_id;references:id"`
	Image Image `gorm:"foreignKey:image_id;references:id"`
}

func (l *Like) TableName() string {
	return LikeTableName
}

type LikeList []Like
