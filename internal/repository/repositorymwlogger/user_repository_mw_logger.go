package repositorymwlogger

import (
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ repository.UserRepository = &UserRepositoryImpl{}

type UserRepositoryImpl struct {
	logger *logrus.Logger

	next repository.UserRepository
}

func NewUserRepository(logger *logrus.Logger, next repository.UserRepository) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		logger: logger,
		next:   next,
	}
}

func (r *UserRepositoryImpl) FindByToken(db *gorm.DB, user *entity.User, token string) error {
	err := r.next.FindByToken(db, user, token)

	fields := logrus.Fields{
		"user":  user,
		"token": token,
	}
	helper.Log(r.logger, fields, err)

	return err
}

func (r *UserRepositoryImpl) CountById(db *gorm.DB, id string) (int64, error) {
	total, err := r.next.CountById(db, id)

	fields := logrus.Fields{
		"id":    id,
		"total": total,
	}
	helper.Log(r.logger, fields, err)

	return total, err
}

func (r *UserRepositoryImpl) Create(db *gorm.DB, entity *entity.User) error {
	err := r.next.Create(db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}

func (r *UserRepositoryImpl) FindById(db *gorm.DB, entity *entity.User, id string) error {
	err := r.next.FindById(db, entity, id)

	fields := logrus.Fields{
		"id":     id,
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}

func (r *UserRepositoryImpl) Update(db *gorm.DB, entity *entity.User) error {
	err := r.next.Update(db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(r.logger, fields, err)

	return err
}
