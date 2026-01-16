package model

type UploadImageRequest struct {
}

type LikeImageRequest struct {
	ImageID int64 `json:"image_id"`
}

type CommentImageRequest struct {
	ImageID int64  `json:"image_id"`
	Comment string `json:"comment"`
}
