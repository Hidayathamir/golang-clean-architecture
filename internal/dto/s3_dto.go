package dto

import (
	"io"
	"time"
)

type S3DeleteObjectRequest struct {
	Bucket string
	Key    string
}

type S3DeleteObjectResponse struct {
	Deleted bool
}

type S3UploadImageRequest struct {
	Key  string
	Body io.Reader
}

type S3PresignUploadRequest struct {
	Key         string
	ContentType string
	ExpiresIn   time.Duration
}

type S3PresignUploadResponse struct {
	UploadURL string
	Key       string
}

type S3ObjectExistsRequest struct {
	Key string
}

type S3ObjectExistsResponse struct {
	Exists bool
}
