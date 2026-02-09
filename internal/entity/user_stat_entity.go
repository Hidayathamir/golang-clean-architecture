package entity

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"gorm.io/gorm"
)

type UserStat struct {
	ID             int64          `gorm:"column:id;primaryKey"`
	UserID         int64          `gorm:"column:user_id"`
	FollowerCount  int            `gorm:"column:follower_count"`
	FollowingCount int            `gorm:"column:following_count"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (u *UserStat) TableName() string {
	return table.UserStat
}
