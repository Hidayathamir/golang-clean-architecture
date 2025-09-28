package repository

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ UserRepository = &UserRepositoryMwLogger{}

type UserRepositoryMwLogger struct {
	Next UserRepository
}

func NewUserRepositoryMwLogger(next UserRepository) *UserRepositoryMwLogger {
	return &UserRepositoryMwLogger{
		Next: next,
	}
}

func (r *UserRepositoryMwLogger) FindByToken(ctx context.Context, db *gorm.DB, user *entity.User, token string) error {
	err := r.Next.FindByToken(ctx, db, user, token)

	fields := logrus.Fields{
		"user":  user,
		"token": token,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) CountById(ctx context.Context, db *gorm.DB, id string) (int64, error) {
	total, err := r.Next.CountById(ctx, db, id)

	fields := logrus.Fields{
		"id":    id,
		"total": total,
	}
	helper.Log(ctx, fields, err)

	return total, err
}

func (r *UserRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := r.Next.Create(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) FindById(ctx context.Context, db *gorm.DB, entity *entity.User, id string) error {
	err := r.Next.FindById(ctx, db, entity, id)

	fields := logrus.Fields{
		"id":     id,
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := r.Next.Update(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	helper.Log(ctx, fields, err)

	return err
}
