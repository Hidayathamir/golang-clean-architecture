package storage

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ S3Client = &S3ClientMwLogger{}

type S3ClientMwLogger struct {
	Next S3Client
}

func NewS3ClientMwLogger(next S3Client) *S3ClientMwLogger {
	return &S3ClientMwLogger{
		Next: next,
	}
}

func (c *S3ClientMwLogger) UploadImage(ctx context.Context, req model.S3UploadImageRequest) (url string, err error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	url, err = c.Next.UploadImage(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"key": req.Key,
		"url": url,
	}
	x.LogMw(ctx, fields, err)

	return url, err
}

func (c *S3ClientMwLogger) Download(ctx context.Context, req model.S3DownloadRequest) (model.S3DownloadResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := c.Next.Download(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"bucket": req.Bucket,
		"key":    req.Key,
		"data":   res.Data,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (c *S3ClientMwLogger) DeleteObject(ctx context.Context, req model.S3DeleteObjectRequest) (model.S3DeleteObjectResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := c.Next.DeleteObject(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"bucket":  req.Bucket,
		"key":     req.Key,
		"deleted": res.Deleted,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}
