package entity

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"gorm.io/gorm"
)

type Image struct {
	ID           int64          `gorm:"column:id;primaryKey"`
	UserID       int64          `gorm:"column:user_id"`
	Caption      string         `gorm:"column:caption"`
	URL          string         `gorm:"column:url"`
	LikeCount    int            `gorm:"column:like_count"`
	CommentCount int            `gorm:"column:comment_count"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (i *Image) TableName() string {
	return table.Image
}

type ImageList []Image
