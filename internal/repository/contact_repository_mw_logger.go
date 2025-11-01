package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
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

func (r *ContactRepositoryMwLogger) FindByIDAndUserID(ctx context.Context, db *gorm.DB, contact *entity.Contact, id string, userID string) error {
	err := r.Next.FindByIDAndUserID(ctx, db, contact, id, userID)

	fields := logrus.Fields{
		"contact": contact,
		"id":      id,
		"userID":  userID,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *ContactRepositoryMwLogger) Search(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) (entity.ContactList, int64, error) {
	contacts, total, err := r.Next.Search(ctx, db, req)

	fields := logrus.Fields{
		"req":      req,
		"contacts": contacts,
		"total":    total,
	}
	logging.Log(ctx, fields, err)

	return contacts, total, err
}

func (r *ContactRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := r.Next.Create(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *ContactRepositoryMwLogger) Delete(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := r.Next.Delete(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *ContactRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := r.Next.Update(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	logging.Log(ctx, fields, err)

	return err
}
