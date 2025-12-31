package todo_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTodoUsecaseImpl_Create_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	TodoProducer := &mock.TodoProducerMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
		TodoProducer:   TodoProducer,
	}

	// ------------------------------------------------------- //

	req := &model.CreateTodoRequest{
		UserID:      testUserID,
		Title:       "Test Todo",
		Description: "Test Description",
	}

	fixedTime := time.UnixMilli(1699432800000).UTC()
	var capturedTodo *entity.Todo
	TodoRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
		capturedTodo = todo
		todo.ID = 1 // simulate auto-generated ID
		todo.CreatedAt = fixedTime
		todo.UpdatedAt = fixedTime
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotZero(t, capturedTodo.ID) // Verify ID was generated
	assert.Equal(t, req.UserID, capturedTodo.UserID)
	assert.Equal(t, req.Title, capturedTodo.Title)
	assert.Equal(t, req.Description, capturedTodo.Description)
	assert.False(t, capturedTodo.IsCompleted)
	assert.False(t, capturedTodo.CompletedAt.Valid)

	assert.Equal(t, capturedTodo.ID, res.ID)
	assert.Equal(t, capturedTodo.Title, res.Title)
	assert.Equal(t, capturedTodo.Description, res.Description)
	assert.Equal(t, capturedTodo.IsCompleted, res.IsCompleted)
	assert.Equal(t, capturedTodo.CompletedAt, res.CompletedAt)
	assert.Equal(t, capturedTodo.CreatedAt, res.CreatedAt)
	assert.Equal(t, capturedTodo.UpdatedAt, res.UpdatedAt)
	assert.Nil(t, err)
}

func TestTodoUsecaseImpl_Create_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.CreateTodoRequest{
		// Missing required fields UserID and Title
		Description: "Test Description",
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse
	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Field validation") // errkit wraps validation errors
}

func TestTodoUsecaseImpl_Create_Fail_Create(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.CreateTodoRequest{
		UserID:      testUserID,
		Title:       "Test Todo",
		Description: "Test Description",
	}

	TodoRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}
