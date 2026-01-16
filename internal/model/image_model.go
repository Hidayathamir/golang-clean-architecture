package model

import "mime/multipart"

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
