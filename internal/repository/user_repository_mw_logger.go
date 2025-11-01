package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
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
	logging.Log(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) CountByID(ctx context.Context, db *gorm.DB, id string) (int64, error) {
	total, err := r.Next.CountByID(ctx, db, id)

	fields := logrus.Fields{
		"id":    id,
		"total": total,
	}
	logging.Log(ctx, fields, err)

	return total, err
}

func (r *UserRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := r.Next.Create(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) FindByID(ctx context.Context, db *gorm.DB, entity *entity.User, id string) error {
	err := r.Next.FindByID(ctx, db, entity, id)

	fields := logrus.Fields{
		"id":     id,
		"entity": entity,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *UserRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, entity *entity.User) error {
	err := r.Next.Update(ctx, db, entity)

	fields := logrus.Fields{
		"entity": entity,
	}
	logging.Log(ctx, fields, err)

	return err
}
