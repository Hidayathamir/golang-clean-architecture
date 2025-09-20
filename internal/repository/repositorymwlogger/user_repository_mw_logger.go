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

	RepositoryImpl[entity.User]
	next repository.UserRepository
}

func NewUserRepository(logger *logrus.Logger, next repository.UserRepository) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		RepositoryImpl: RepositoryImpl[entity.User]{
			logger: logger,
			next:   next,
		},
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
