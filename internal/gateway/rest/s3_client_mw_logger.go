package rest

import (
	"context"

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

func (c *S3ClientMwLogger) Download(ctx context.Context, bucket, key string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	path, err := c.Next.Download(ctx, bucket, key)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"bucket": bucket,
		"key":    key,
		"path":   path,
	}
	x.LogMw(ctx, fields, err)

	return path, err
}

func (c *S3ClientMwLogger) DeleteObject(ctx context.Context, bucket, key string) (bool, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	ok, err := c.Next.DeleteObject(ctx, bucket, key)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"bucket": bucket,
		"key":    key,
		"ok":     ok,
	}
	x.LogMw(ctx, fields, err)

	return ok, err
}
