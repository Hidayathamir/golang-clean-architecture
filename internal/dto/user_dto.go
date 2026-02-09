package dto

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VerifyUserRequest struct {
	Token string `validate:"required"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Name     string `json:"name"     validate:"required,max=100"`
}

type UpdateUserRequest struct {
	ID       int64  `json:"-"        validate:"required"`
	Password string `json:"password" validate:"max=100"`
	Name     string `json:"name"     validate:"max=100"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type UserLoginResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LogoutUserRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type GetUserRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type FollowUserRequest struct {
	FollowingID int64 `json:"following_id"`
}

type NotifyUserBeingFollowedRequest struct {
	FollowerID  int64
	FollowingID int64
}

type BatchUpdateUserFollowStatsRequest struct {
	UserIncreaseFollowerFollowingCountList UserIncreaseFollowerFollowingCountList
}

type UserIncreaseFollowerFollowingCount struct {
	UserID         int64
	FollowerCount  int
	FollowingCount int
}

func (u UserIncreaseFollowerFollowingCount) HasFollowerCount() bool {
	return u.FollowerCount > 0
}

func (u UserIncreaseFollowerFollowingCount) HasFollowingCount() bool {
	return u.FollowingCount > 0
}

func (u UserIncreaseFollowerFollowingCount) HasFollowerCountAndFollowingCount() bool {
	return u.HasFollowerCount() && u.HasFollowingCount()
}

type UserIncreaseFollowerFollowingCountList []UserIncreaseFollowerFollowingCount

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
