package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
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
	ok, err := c.Next.Download(ctx, bucket, key)

	fields := logrus.Fields{
		"bucket": bucket,
		"key":    key,
		"ok":     ok,
	}
	l.LogMw(ctx, fields, err)

	return ok, err
}

func (c *S3ClientMwLogger) DeleteObject(ctx context.Context, bucket, key string) (bool, error) {
	ok, err := c.Next.DeleteObject(ctx, bucket, key)

	fields := logrus.Fields{
		"bucket": bucket,
		"key":    key,
		"ok":     ok,
	}
	l.LogMw(ctx, fields, err)

	return ok, err
}
