package todo

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *TodoUsecaseImpl) List(ctx context.Context, req *model.ListTodoRequest) (model.TodoResponseList, int64, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, 0, errkit.AddFuncName("todo.(*TodoUsecaseImpl).List", err)
	}

	todos, total, err := u.TodoRepository.List(ctx, u.DB.WithContext(ctx), req)
	if err != nil {
		return nil, 0, errkit.AddFuncName("todo.(*TodoUsecaseImpl).List", err)
	}

	res := make(model.TodoResponseList, len(todos))
	converter.EntityTodoListToModelTodoResponseList(todos, res)

	return res, total, nil
}
