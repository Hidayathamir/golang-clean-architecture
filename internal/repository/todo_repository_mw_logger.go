package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ TodoRepository = &TodoRepositoryMwLogger{}

type TodoRepositoryMwLogger struct {
	Next TodoRepository
}

func NewTodoRepositoryMwLogger(next TodoRepository) *TodoRepositoryMwLogger {
	return &TodoRepositoryMwLogger{
		Next: next,
	}
}

func (r *TodoRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	err := r.Next.Create(ctx, db, todo)

	fields := logrus.Fields{
		"todo": todo,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *TodoRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	err := r.Next.Update(ctx, db, todo)

	fields := logrus.Fields{
		"todo": todo,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *TodoRepositoryMwLogger) Delete(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	err := r.Next.Delete(ctx, db, todo)

	fields := logrus.Fields{
		"todo": todo,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *TodoRepositoryMwLogger) FindByIDAndUserID(ctx context.Context, db *gorm.DB, todo *entity.Todo, id string, userID string) error {
	err := r.Next.FindByIDAndUserID(ctx, db, todo, id, userID)

	fields := logrus.Fields{
		"id":      id,
		"user_id": userID,
		"todo":    todo,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (r *TodoRepositoryMwLogger) List(ctx context.Context, db *gorm.DB, req *model.ListTodoRequest) (entity.TodoList, int64, error) {
	todos, total, err := r.Next.List(ctx, db, req)

	fields := logrus.Fields{
		"req":   req,
		"todos": todos,
		"total": total,
	}
	logging.Log(ctx, fields, err)

	return todos, total, err
}
