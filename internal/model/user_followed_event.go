package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type UserFollowedEvent struct {
	ID          int64          `json:"id"`
	FollowerID  int64          `json:"follower_id"`
	FollowingID int64          `json:"following_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (u *UserFollowedEvent) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}
