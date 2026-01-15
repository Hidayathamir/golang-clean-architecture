package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	UserTableName            = "users"
	UserColumnID             = "id"
	UserColumnUsername       = "username"
	UserColumnPassword       = "password"
	UserColumnName           = "name"
	UserColumnFollowerCount  = "follower_count"
	UserColumnFollowingCount = "following_count"
	UserColumnCreatedAt      = "created_at"
	UserColumnUpdatedAt      = "updated_at"
	UserColumnDeletedAt      = "deleted_at"
)

type User struct {
	ID             int64          `gorm:"column:id;primaryKey"`
	Username       string         `gorm:"column:username;uniqueIndex"`
	Password       string         `gorm:"column:password"`
	Name           string         `gorm:"column:name"`
	FollowerCount  int            `gorm:"column:follower_count"`
	FollowingCount int            `gorm:"column:following_count"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (u *User) TableName() string {
	return UserTableName
}

type UserList []User
