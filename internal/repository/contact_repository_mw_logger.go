package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ ContactRepository = &ContactRepositoryMwLogger{}

type ContactRepositoryMwLogger struct {
	Next ContactRepository
}

func NewContactRepositoryMwLogger(next ContactRepository) *ContactRepositoryMwLogger {
	return &ContactRepositoryMwLogger{
		Next: next,
	}
}

func (r *ContactRepositoryMwLogger) FindByIDAndUserID(ctx context.Context, db *gorm.DB, contact *entity.Contact, id int64, userID int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByIDAndUserID(ctx, db, contact, id, userID)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"contact": contact,
		"id":      id,
		"userID":  userID,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *ContactRepositoryMwLogger) Search(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) (entity.ContactList, int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	contacts, total, err := r.Next.Search(ctx, db, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req":      req,
		"contacts": contacts,
		"total":    total,
	}
	x.LogMw(ctx, fields, err)

	return contacts, total, err
}

func (r *ContactRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
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

func (r *ContactRepositoryMwLogger) Delete(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Delete(ctx, db, entity)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"entity": entity,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *ContactRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Update(ctx, db, entity)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"entity": entity,
	}
	x.LogMw(ctx, fields, err)

	return err
}
