package todo

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ TodoUsecase = &TodoUsecaseMwTelemetry{}

type TodoUsecaseMwTelemetry struct {
	Next TodoUsecase
}

func NewTodoUsecaseMwTelemetry(next TodoUsecase) *TodoUsecaseMwTelemetry {
	return &TodoUsecaseMwTelemetry{
		Next: next,
	}
}

func (u *TodoUsecaseMwTelemetry) Create(ctx context.Context, req *model.CreateTodoRequest) (*model.TodoResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Create(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *TodoUsecaseMwTelemetry) Get(ctx context.Context, req *model.GetTodoRequest) (*model.TodoResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Get(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *TodoUsecaseMwTelemetry) List(ctx context.Context, req *model.ListTodoRequest) (model.TodoResponseList, int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	list, total, err := u.Next.List(ctx, req)
	telemetry.RecordError(span, err)

	return list, total, err
}

func (u *TodoUsecaseMwTelemetry) Update(ctx context.Context, req *model.UpdateTodoRequest) (*model.TodoResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Update(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *TodoUsecaseMwTelemetry) Delete(ctx context.Context, req *model.DeleteTodoRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Delete(ctx, req)
	telemetry.RecordError(span, err)

	return err
}

func (u *TodoUsecaseMwTelemetry) Complete(ctx context.Context, req *model.CompleteTodoRequest) (*model.TodoResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Complete(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}
