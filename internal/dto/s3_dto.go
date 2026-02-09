package dto

import "io"

type S3DownloadRequest struct {
	Bucket string
	Key    string
}

type S3DownloadResponse struct {
	Data string
}

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
