package model

import (
	"time"

	"gorm.io/gorm"
)

type ImageUploadedEvent struct {
	ID           int64          `json:"id"`
	UserID       int64          `json:"user_id"`
	URL          string         `json:"url"`
	LikeCount    int            `json:"like_count"`
	CommentCount int            `json:"comment_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}
