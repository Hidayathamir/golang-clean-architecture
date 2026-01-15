package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	ImageTableName          = "images"
	ImageColumnID           = "id"
	ImageColumnUserID       = "user_id"
	ImageColumnURL          = "url"
	ImageColumnLikeCount    = "like_count"
	ImageColumnCommentCount = "comment_count"
	ImageColumnCreatedAt    = "created_at"
	ImageColumnUpdatedAt    = "updated_at"
	ImageColumnDeletedAt    = "deleted_at"
)

type Image struct {
	ID           int64          `gorm:"column:id;primaryKey"`
	UserID       int64          `gorm:"column:user_id"`
	URL          string         `gorm:"column:url"`
	LikeCount    int            `gorm:"column:like_count"`
	CommentCount int            `gorm:"column:comment_count"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`

	User User `gorm:"foreignKey:user_id;references:id"`
}

func (i *Image) TableName() string {
	return ImageTableName
}

type ImageList []Image
