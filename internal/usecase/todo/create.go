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

	todo := new(entity.Todo)
	converter.ModelCreateTodoRequestToEntityTodo(req, todo, uuid.New().String())

	if err := u.TodoRepository.Create(ctx, u.DB.WithContext(ctx), todo); err != nil {
		return nil, errkit.AddFuncName("todo.(*TodoUsecaseImpl).Create", err)
	}

	res := new(model.TodoResponse)
	converter.EntityTodoToModelTodoResponse(todo, res)

	return res, nil
}
