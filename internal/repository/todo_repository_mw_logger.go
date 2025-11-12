package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
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
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Create(ctx, db, todo)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"todo": todo,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *TodoRepositoryMwLogger) Update(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Update(ctx, db, todo)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"todo": todo,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *TodoRepositoryMwLogger) Delete(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Delete(ctx, db, todo)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"todo": todo,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *TodoRepositoryMwLogger) FindByIDAndUserID(ctx context.Context, db *gorm.DB, todo *entity.Todo, id string, userID int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByIDAndUserID(ctx, db, todo, id, userID)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":      id,
		"user_id": userID,
		"todo":    todo,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (r *TodoRepositoryMwLogger) List(ctx context.Context, db *gorm.DB, req *model.ListTodoRequest) (entity.TodoList, int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	todos, total, err := r.Next.List(ctx, db, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req":   req,
		"todos": todos,
		"total": total,
	}
	x.LogMw(ctx, fields, err)

	return todos, total, err
}
