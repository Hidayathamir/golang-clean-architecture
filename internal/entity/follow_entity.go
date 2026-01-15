package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	FollowTableName         = "follows"
	FollowColumnID          = "id"
	FollowColumnFollowerID  = "follower_id"
	FollowColumnFollowingID = "following_id"
	FollowColumnCreatedAt   = "created_at"
	FollowColumnUpdatedAt   = "updated_at"
	FollowColumnDeletedAt   = "deleted_at"
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
	return FollowTableName
}

type FollowList []Follow
