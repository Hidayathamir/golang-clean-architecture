package todo

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *TodoUsecaseImpl) Complete(ctx context.Context, req *model.CompleteTodoRequest) (*model.TodoResponse, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("todo.(*TodoUsecaseImpl).Complete", err)
	}

	todo := new(entity.Todo)
	if err := u.TodoRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), todo, req.ID, req.UserID); err != nil {
		return nil, errkit.AddFuncName("todo.(*TodoUsecaseImpl).Complete", err)
	}

	if !todo.IsCompleted {
		now := time.Now().UnixMilli()
		todo.IsCompleted = true
		todo.CompletedAt = &now

		if err := u.TodoRepository.Update(ctx, u.DB.WithContext(ctx), todo); err != nil {
			return nil, errkit.AddFuncName("todo.(*TodoUsecaseImpl).Complete", err)
		}

		event := new(model.TodoCompletedEvent)
		converter.EntityTodoToModelTodoCompletedEvent(todo, event)
		if err := u.TodoProducer.Send(ctx, event); err != nil {
			return nil, errkit.AddFuncName("todo.(*TodoUsecaseImpl).Complete", err)
		}
	}

	res := new(model.TodoResponse)
	converter.EntityTodoToModelTodoResponse(todo, res)

	return res, nil
}
