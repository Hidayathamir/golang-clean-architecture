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

// DeleteImageRequest is the input for DELETE /api/images/:imageId. The caller
// (owner check) comes from the auth context.
type DeleteImageRequest struct {
	ID int64 `json:"-" validate:"required"`
}

// UpdateImageRequest is the input for PATCH /api/images/:imageId. The caller
// (owner check) comes from the auth context.
type UpdateImageRequest struct {
	ID      int64  `json:"-"       validate:"required"`
	Caption string `json:"caption" validate:"required,max=2200"`
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

// FeedRequest is the input for GET /api/feed. The user is taken from auth
// context; Page carries the cursor pagination input.
type FeedRequest struct {
	Page PageRequest
}

type ImageResponseList []ImageResponse

// FeedResponse is the feed page: images from followed users, newest first.
type FeedResponse = CursorPage[ImageResponseList]

// UserImagesRequest is the input for GET /api/users/:id/images. The target user
// comes from the path; Page carries the cursor pagination input.
type UserImagesRequest struct {
	UserID int64
	Page   PageRequest
}

// UserImagesResponse is the user's post grid page, newest first.
type UserImagesResponse = CursorPage[ImageResponseList]

// ExploreRequest is the input for GET /api/explore. Page carries the cursor
// pagination input; discovery order is newest first.
type ExploreRequest struct {
	Page PageRequest
}

// ExploreResponse is a page of discovered posts, newest first.
type ExploreResponse = CursorPage[ImageResponseList]

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
