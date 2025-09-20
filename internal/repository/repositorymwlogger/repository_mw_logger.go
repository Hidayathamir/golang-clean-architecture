package repositorymwlogger

import (
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ repository.Repository[any] = &RepositoryImpl[any]{}

type RepositoryImpl[T any] struct {
	logger *logrus.Logger

	next repository.Repository[T]
}

func NewRepository(logger *logrus.Logger, next repository.Repository[any]) *RepositoryImpl[any] {
	return &RepositoryImpl[any]{
		logger: logger,
		next:   next,
	}
}

func (r *RepositoryImpl[T]) CountById(db *gorm.DB, id any) (int64, error) {
	total, err := r.next.CountById(db, id)

	fields := logrus.Fields{
		"id":    id,
		"total": total,
	}
	helper.Log(r.logger, fields, err)

	return total, err
}

func (r *RepositoryImpl[T]) Create(db *gorm.DB, entity *T) error {
	err := r.next.Create(db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}

func (r *RepositoryImpl[T]) Delete(db *gorm.DB, entity *T) error {
	err := r.next.Delete(db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}

func (r *RepositoryImpl[T]) FindById(db *gorm.DB, entity *T, id any) error {
	err := r.next.FindById(db, entity, id)

	fields := logrus.Fields{
		"id":     id,
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}

func (r *RepositoryImpl[T]) Update(db *gorm.DB, entity *T) error {
	err := r.next.Update(db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}
