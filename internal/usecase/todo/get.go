package todo

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *TodoUsecaseImpl) Get(ctx context.Context, req *model.GetTodoRequest) (*model.TodoResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("todo.(*TodoUsecaseImpl).Get", err)
	}

	todo := new(entity.Todo)
	if err := u.TodoRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), todo, req.ID, req.UserID); err != nil {
		return nil, errkit.AddFuncName("todo.(*TodoUsecaseImpl).Get", err)
	}

	return converter.TodoToResponse(todo), nil
}
