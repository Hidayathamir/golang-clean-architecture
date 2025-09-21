package repositorymwlogger

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ repository.ContactRepository = &ContactRepositoryImpl{}

type ContactRepositoryImpl struct {
	logger *logrus.Logger

	next repository.ContactRepository
}

func NewContactRepository(logger *logrus.Logger, next repository.ContactRepository) *ContactRepositoryImpl {
	return &ContactRepositoryImpl{
		logger: logger,
		next:   next,
	}
}

func (r *ContactRepositoryImpl) FindByIdAndUserId(ctx context.Context, db *gorm.DB, contact *entity.Contact, id string, userId string) error {
	err := r.next.FindByIdAndUserId(ctx, db, contact, id, userId)

	fields := logrus.Fields{
		"contact": contact,
		"id":      id,
		"userId":  userId,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *ContactRepositoryImpl) Search(ctx context.Context, db *gorm.DB, req *model.SearchContactRequest) ([]entity.Contact, int64, error) {
	contacts, total, err := r.next.Search(ctx, db, req)

	fields := logrus.Fields{
		"req":      req,
		"contacts": contacts,
		"total":    total,
	}
	helper.Log(ctx, fields, err)

	return contacts, total, err
}

func (r *ContactRepositoryImpl) Create(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := r.next.Create(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *ContactRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := r.next.Delete(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *ContactRepositoryImpl) Update(ctx context.Context, db *gorm.DB, entity *entity.Contact) error {
	err := r.next.Update(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}
