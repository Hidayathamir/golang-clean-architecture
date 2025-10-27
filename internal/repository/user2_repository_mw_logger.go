package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ User2Repository = &User2RepositoryMwLogger{}

type User2RepositoryMwLogger struct {
	Next User2Repository
}

func NewUser2RepositoryMwLogger(next User2Repository) *User2RepositoryMwLogger {
	return &User2RepositoryMwLogger{
		Next: next,
	}
}

func (r *User2RepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, user *entity.User2) error {
	err := r.Next.Create(ctx, db, user)

	fields := logrus.Fields{
		"user": user,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *User2RepositoryMwLogger) FindByEmail(ctx context.Context, db *gorm.DB, user *entity.User2, email string) error {
	err := r.Next.FindByEmail(ctx, db, user, email)

	fields := logrus.Fields{
		"email": email,
		"user":  user,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *User2RepositoryMwLogger) FindByID(ctx context.Context, db *gorm.DB, user *entity.User2, id string) error {
	err := r.Next.FindByID(ctx, db, user, id)

	fields := logrus.Fields{
		"id":   id,
		"user": user,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *User2RepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, user *entity.User2) error {
	err := r.Next.Update(ctx, db, user)

	fields := logrus.Fields{
		"user": user,
	}
	logging.Log(ctx, fields, err)

	return err
}
