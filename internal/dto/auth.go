package dto

import (
	"time"

	"gorm.io/gorm"
)

type UserAuth struct {
	ID             int64
	Username       string
	Name           string
	FollowerCount  int
	FollowingCount int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}
