package entity

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"gorm.io/gorm"
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
	return table.Like
}

type LikeList []Like
