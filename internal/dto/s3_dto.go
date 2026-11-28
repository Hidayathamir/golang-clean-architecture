package dto

import "io"

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
