package entity

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"gorm.io/gorm"
)

type Follow struct {
	ID          int64          `gorm:"column:id;primaryKey"`
	FollowerID  int64          `gorm:"column:follower_id"`
	FollowingID int64          `gorm:"column:following_id"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`

	Follower  User `gorm:"foreignKey:follower_id;references:id"`
	Following User `gorm:"foreignKey:following_id;references:id"`
}

func (f *Follow) TableName() string {
	return table.Follow
}

type FollowList []Follow
