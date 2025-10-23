package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
)

func TodoToResponse(todo *entity.Todo) *model.TodoResponse {
	return &model.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
		CompletedAt: todo.CompletedAt,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

func TodosToResponses(todos []entity.Todo) []model.TodoResponse {
	responses := make([]model.TodoResponse, 0, len(todos))
	for _, todo := range todos {
		responses = append(responses, *TodoToResponse(&todo))
	}
	return responses
}

func TodoToCompletedEvent(todo *entity.Todo) *model.TodoCompletedEvent {
	return &model.TodoCompletedEvent{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		CompletedAt: todo.CompletedAt,
	}
}
