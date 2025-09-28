package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/helper"

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

func (r *ContactRepositoryMwLogger) FindByIdAndUserId(ctx context.Context, db *gorm.DB, contact *entity.Contact, id string, userId string) error {
	err := r.Next.FindByIdAndUserId(ctx, db, contact, id, userId)

	fields := logrus.Fields{
		"contact": contact,
		"id":      id,
		"userId":  userId,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *ContactRepositoryMwLogger) Search(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) ([]entity.Contact, int64, error) {
	contacts, total, err := r.Next.Search(ctx, db, req)

	fields := logrus.Fields{
		"req":      req,
		"contacts": contacts,
		"total":    total,
	}
	helper.Log(ctx, fields, err)

	return contacts, total, err
}

func (r *ContactRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := r.Next.Create(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *ContactRepositoryMwLogger) Delete(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := r.Next.Delete(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *ContactRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := r.Next.Update(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}
