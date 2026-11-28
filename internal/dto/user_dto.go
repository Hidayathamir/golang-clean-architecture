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

// ProfileResponse is the public profile for GET /api/users/:id.
type ProfileResponse struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	FollowerCount  int    `json:"follower_count"`
	FollowingCount int    `json:"following_count"`
}

// GetProfileRequest is the input for GET /api/users/:id.
type GetProfileRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type UserResponseList []UserResponse

// FollowersRequest is the input for GET /api/users/:id/followers. The target
// user comes from the path; Page carries the cursor pagination input.
type FollowersRequest struct {
	UserID int64
	Page   PageRequest
}

// FollowersResponse is a page of the target user's followers.
type FollowersResponse = CursorPage[UserResponseList]

// FollowingsRequest is the input for GET /api/users/:id/following. The target
// user comes from the path; Page carries the cursor pagination input.
type FollowingsRequest struct {
	UserID int64
	Page   PageRequest
}

// FollowingsResponse is a page of the users the target user follows.
type FollowingsResponse = CursorPage[UserResponseList]

type FollowUserRequest struct {
	FollowingID int64 `json:"following_id"`
}

// UnfollowUserRequest is the input for DELETE /api/users/:id/follow. The
// follower comes from the auth context.
type UnfollowUserRequest struct {
	FollowingID int64 `json:"-" validate:"required"`
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
	return u.FollowerCount != 0
}

func (u UserIncreaseFollowerFollowingCount) HasFollowingCount() bool {
	return u.FollowingCount != 0
}

func (u UserIncreaseFollowerFollowingCount) HasFollowerCountAndFollowingCount() bool {
	return u.HasFollowerCount() && u.HasFollowingCount()
}

type UserIncreaseFollowerFollowingCountList []UserIncreaseFollowerFollowingCount

type UserFollowedEvent struct {
	ID          int64          `json:"id"`
	FollowerID  int64          `json:"follower_id"`
	FollowingID int64          `json:"following_id"`
	Count       int            `json:"count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (u *UserFollowedEvent) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}
