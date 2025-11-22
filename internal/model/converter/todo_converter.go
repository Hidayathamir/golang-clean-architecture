package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
)

func ModelCreateTodoRequestToEntityTodo(req *model.CreateTodoRequest, todo *entity.Todo) {
	if req == nil || todo == nil {
		return
	}

	todo.UserID = req.UserID
	todo.Title = req.Title
	todo.Description = req.Description
	todo.IsCompleted = false
	todo.CompletedAt = nil
}

func ModelUpdateTodoRequestToEntityTodo(req *model.UpdateTodoRequest, todo *entity.Todo) {
	if req == nil || todo == nil {
		return
	}

	todo.Title = req.Title
	todo.Description = req.Description
}

func EntityTodoToModelTodoResponse(todo *entity.Todo, res *model.TodoResponse) {
	if todo == nil || res == nil {
		return
	}

	res.ID = todo.ID
	res.Title = todo.Title
	res.Description = todo.Description
	res.IsCompleted = todo.IsCompleted
	res.CompletedAt = todo.CompletedAt
	res.CreatedAt = todo.CreatedAt
	res.UpdatedAt = todo.UpdatedAt
}

func EntityTodoListToModelTodoResponseList(todos entity.TodoList, res model.TodoResponseList) {
	if res == nil {
		return
	}

	for i := range todos {
		EntityTodoToModelTodoResponse(&todos[i], &res[i])
	}
}

func EntityTodoToModelTodoCompletedEvent(todo *entity.Todo, event *model.TodoCompletedEvent) {
	if todo == nil || event == nil {
		return
	}

	event.ID = todo.ID
	event.UserID = todo.UserID
	event.Title = todo.Title
	event.Description = todo.Description
	event.CompletedAt = todo.CompletedAt
}
