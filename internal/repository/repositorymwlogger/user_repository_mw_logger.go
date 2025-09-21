package repositorymwlogger

import (
	"context"
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

func (r *UserRepositoryImpl) FindByToken(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
	err := r.next.FindByToken(ctx, db, user, token)

	fields := logrus.Fields{
		"user":  user,
		"token": token,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *UserRepositoryImpl) CountById(ctx context.Context, db *gorm.DB, id string) (int64, error) {
	total, err := r.next.CountById(ctx, db, id)

	fields := logrus.Fields{
		"id":    id,
		"total": total,
	}
	helper.Log(ctx, fields, err)

	return total, err
}

func (r *UserRepositoryImpl) Create(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := r.next.Create(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *UserRepositoryImpl) FindById(ctx context.Context, db *gorm.DB, entity *entity.User, id string) error {
	err := r.next.FindById(ctx, db, entity, id)

	fields := logrus.Fields{
		"id":     id,
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := r.next.Update(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}
