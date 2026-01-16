package model

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type ImageResponse struct {
	ID           int64          `json:"id"`
	UserID       int64          `json:"user_id"`
	URL          string         `json:"url"`
	LikeCount    int            `json:"like_count"`
	CommentCount int            `json:"comment_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type UploadImageRequest struct {
	File *multipart.FileHeader `validate:"required"`
}

type LikeImageRequest struct {
	ImageID int64 `json:"image_id"`
}

type CommentImageRequest struct {
	ImageID int64  `json:"image_id"`
	Comment string `json:"comment"`
}

type GetImageRequest struct {
	ID int64
}
