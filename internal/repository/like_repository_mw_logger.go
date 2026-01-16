package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ LikeRepository = &LikeRepositoryMwLogger{}

type LikeRepositoryMwLogger struct {
	Next LikeRepository
}

func NewLikeRepositoryMwLogger(next LikeRepository) *LikeRepositoryMwLogger {
	return &LikeRepositoryMwLogger{
		Next: next,
	}
}

func (r *LikeRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.Like) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Create(ctx, db, entity)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"entity": entity,
	}
	x.LogMw(ctx, fields, err)

	return err
}
