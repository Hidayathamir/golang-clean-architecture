package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ CommentRepository = &CommentRepositoryMwLogger{}

type CommentRepositoryMwLogger struct {
	Next CommentRepository
}

func NewCommentRepositoryMwLogger(next CommentRepository) *CommentRepositoryMwLogger {
	return &CommentRepositoryMwLogger{
		Next: next,
	}
}

func (r *CommentRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, comment *entity.Comment) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.Create(ctx, db, comment)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"comment": comment,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *CommentRepositoryMwLogger) FindByID(ctx context.Context, db *gorm.DB, comment *entity.Comment, id int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByID(ctx, db, comment, id)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":      id,
		"comment": comment,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *CommentRepositoryMwLogger) DeleteByID(ctx context.Context, db *gorm.DB, comment *entity.Comment) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.DeleteByID(ctx, db, comment)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"comment": comment,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *CommentRepositoryMwLogger) FindByImageID(ctx context.Context, db *gorm.DB, commentList *entity.CommentList, imageID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByImageID(ctx, db, commentList, imageID, cursorCreatedAt, cursorID, limit)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"commentList": commentList,
		"imageID":     imageID,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}
