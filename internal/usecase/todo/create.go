package todo

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/google/uuid"
)

func (u *TodoUsecaseImpl) Create(ctx context.Context, req *model.CreateTodoRequest) (*model.TodoResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("todo.(*TodoUsecaseImpl).Create", err)
	}

	todo := &entity.Todo{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		IsCompleted: false,
	}

	if err := u.TodoRepository.Create(ctx, u.DB.WithContext(ctx), todo); err != nil {
		return nil, errkit.AddFuncName("todo.(*TodoUsecaseImpl).Create", err)
	}

	return converter.TodoToResponse(todo), nil
}
