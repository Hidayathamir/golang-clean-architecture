package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ ImageRepository = &ImageRepositoryMwLogger{}

type ImageRepositoryMwLogger struct {
	Next ImageRepository
}

func NewImageRepositoryMwLogger(next ImageRepository) *ImageRepositoryMwLogger {
	return &ImageRepositoryMwLogger{
		Next: next,
	}
}

func (r *ImageRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.Image) error {
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
