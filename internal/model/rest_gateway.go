package model

import "io"

type PaymentRefundRequest struct {
	TransactionID string
}

type PaymentRefundResponse struct {
	Success bool
}

type PaymentGetStatusRequest struct {
	TransactionID string
}

type PaymentGetStatusResponse struct {
	Status string
}

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

type SlackGetChannelListRequest struct{}

type SlackGetChannelListResponse struct {
	Channels []string
}

type SlackIsConnectedRequest struct{}

type SlackIsConnectedResponse struct {
	Connected bool
}

type S3UploadImageRequest struct {
	Key  string
	Body io.Reader
}
