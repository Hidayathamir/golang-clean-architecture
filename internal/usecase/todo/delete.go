package todo

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *TodoUsecaseImpl) Delete(ctx context.Context, req *model.DeleteTodoRequest) error {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName("todo.(*TodoUsecaseImpl).Delete", err)
	}

	todo := new(entity.Todo)
	if err := u.TodoRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), todo, req.ID, req.UserID); err != nil {
		return errkit.AddFuncName("todo.(*TodoUsecaseImpl).Delete", err)
	}

	if err := u.TodoRepository.Delete(ctx, u.DB.WithContext(ctx), todo); err != nil {
		return errkit.AddFuncName("todo.(*TodoUsecaseImpl).Delete", err)
	}

	return nil
}
