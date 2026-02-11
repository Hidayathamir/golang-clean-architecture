package dto

import (
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

type ImageResponse struct {
	ID           int64          `json:"id"`
	UserID       int64          `json:"user_id"`
	Caption      string         `json:"caption"`
	URL          string         `json:"url"`
	LikeCount    int            `json:"like_count"`
	CommentCount int            `json:"comment_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type UploadImageRequest struct {
	File    *multipart.FileHeader `validate:"required"`
	Caption string
}

type LikeImageRequest struct {
	ImageID int64 `json:"image_id" validate:"required"`
}

type CommentImageRequest struct {
	ImageID int64  `json:"image_id" validate:"required"`
	Comment string `json:"comment"  validate:"required"`
}

type GetImageRequest struct {
	ID int64 `validate:"required"`
}

type LikeResponse struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	ImageID   int64          `json:"image_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type LikeResponseList []LikeResponse

type CommentResponse struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	ImageID   int64          `json:"image_id"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type CommentResponseList []CommentResponse

type GetLikeRequest struct {
	ImageID int64
}

type GetCommentRequest struct {
	ImageID int64
}

type NotifyFollowerOnUploadRequest struct {
	UserID int64
	URL    string
}

type SyncImageToElasticsearchRequest struct {
	ID           int64
	UserID       int64
	Caption      string
	URL          string
	LikeCount    int
	CommentCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

type NotifyUserImageCommentedRequest struct {
	ImageID         int64
	CommenterUserID int64
}

type BatchUpdateImageCommentCountRequest struct {
	ImageIncreaseCommentCountList ImageIncreaseCommentCountList
}

type ImageIncreaseCommentCount struct {
	ImageID int64
	Count   int
}

type ImageIncreaseCommentCountList []ImageIncreaseCommentCount

type NotifyUserImageLikedRequest struct {
	ImageID     int64
	LikerUserID int64
}

type BatchUpdateImageLikeCountRequest struct {
	ImageIncreaseLikeCountList ImageIncreaseLikeCountList
}

type ImageIncreaseLikeCount struct {
	ImageID int64
	Count   int
}

type ImageIncreaseLikeCountList []ImageIncreaseLikeCount

type ImageUploadedEvent struct {
	ID           int64          `json:"id"`
	UserID       int64          `json:"user_id"`
	Caption      string         `json:"caption"`
	URL          string         `json:"url"`
	LikeCount    int            `json:"like_count"`
	CommentCount int            `json:"comment_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type ImageDocument struct {
	ID           int64          `json:"id"`
	UserID       int64          `json:"user_id"`
	Caption      string         `json:"caption"`
	URL          string         `json:"url"`
	LikeCount    int            `json:"like_count"`
	CommentCount int            `json:"comment_count"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

type ImageLikedEvent struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	ImageID   int64          `json:"image_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type ImageCommentedEvent struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	ImageID   int64          `json:"image_id"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
