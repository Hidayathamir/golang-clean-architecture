package model

import (
	"time"

	"gorm.io/gorm"
)

type ImageCommentedEvent struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	ImageID   int64          `json:"image_id"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
