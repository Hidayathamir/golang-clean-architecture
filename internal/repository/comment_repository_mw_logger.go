package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
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

	err := r.Next.Create(ctx, db, comment)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"comment": comment,
	}
	x.LogMw(ctx, fields, err)

	return err
}
