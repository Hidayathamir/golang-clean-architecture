package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"gorm.io/gorm"
)

var _ TodoRepository = &TodoRepositoryMwTelemetry{}

type TodoRepositoryMwTelemetry struct {
	Next TodoRepository
}

func NewTodoRepositoryMwTelemetry(next TodoRepository) *TodoRepositoryMwTelemetry {
	return &TodoRepositoryMwTelemetry{
		Next: next,
	}
}

func (r *TodoRepositoryMwTelemetry) Create(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Create(ctx, db, todo)
	telemetry.RecordError(span, err)

	return err
}

func (r *TodoRepositoryMwTelemetry) Update(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Update(ctx, db, todo)
	telemetry.RecordError(span, err)

	return err
}

func (r *TodoRepositoryMwTelemetry) Delete(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.Delete(ctx, db, todo)
	telemetry.RecordError(span, err)

	return err
}

func (r *TodoRepositoryMwTelemetry) FindByIDAndUserID(ctx context.Context, db *gorm.DB, todo *entity.Todo, id string, userID string) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := r.Next.FindByIDAndUserID(ctx, db, todo, id, userID)
	telemetry.RecordError(span, err)

	return err
}

func (r *TodoRepositoryMwTelemetry) List(ctx context.Context, db *gorm.DB, req *model.ListTodoRequest) (entity.TodoList, int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	list, total, err := r.Next.List(ctx, db, req)
	telemetry.RecordError(span, err)

	return list, total, err
}
