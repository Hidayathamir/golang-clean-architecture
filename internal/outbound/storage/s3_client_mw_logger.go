package storage

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
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

func (c *S3ClientMwLogger) UploadImage(ctx context.Context, req dto.S3UploadImageRequest) (url string, err error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	url, err = retrykit.StorageRetryWithData(ctx, func() (string, error) {
		return c.Next.UploadImage(ctx, req)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"key": req.Key,
		"url": url,
	}
	logkit.LogMw(ctx, fields, err)

	return url, err
}

func (c *S3ClientMwLogger) PresignUploadImage(ctx context.Context, req dto.S3PresignUploadRequest) (dto.S3PresignUploadResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := retrykit.StorageRetryWithData(ctx, func() (dto.S3PresignUploadResponse, error) {
		return c.Next.PresignUploadImage(ctx, req)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"key": req.Key,
	}
	logkit.LogMw(ctx, fields, err)

	return res, err
}

func (c *S3ClientMwLogger) ObjectExists(ctx context.Context, req dto.S3ObjectExistsRequest) (dto.S3ObjectExistsResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := retrykit.StorageRetryWithData(ctx, func() (dto.S3ObjectExistsResponse, error) {
		return c.Next.ObjectExists(ctx, req)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"key":    req.Key,
		"exists": res.Exists,
	}
	logkit.LogMw(ctx, fields, err)

	return res, err
}

func (c *S3ClientMwLogger) DeleteObject(ctx context.Context, req dto.S3DeleteObjectRequest) (dto.S3DeleteObjectResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := retrykit.StorageRetryWithData(ctx, func() (dto.S3DeleteObjectResponse, error) {
		return c.Next.DeleteObject(ctx, req)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"bucket":  req.Bucket,
		"key":     req.Key,
		"deleted": res.Deleted,
	}
	logkit.LogMw(ctx, fields, err)

	return res, err
}
