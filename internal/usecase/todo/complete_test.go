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

func TestTodoUsecaseImpl_Complete_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	TodoProducer := &mock.TodoProducerMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
		TodoProducer:   TodoProducer,
	}

	// ------------------------------------------------------- //

	const todoID int64 = 10
	req := &model.CompleteTodoRequest{
		ID:     todoID,
		UserID: testUserID,
	}

	now := time.Now().UnixMilli()
	incompleteTodo := &entity.Todo{
		ID:          todoID,
		UserID:      req.UserID,
		Title:       "Test Todo",
		Description: "Test Description",
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	var updatedTodo *entity.Todo
	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id int64, userID int64) error {
		*todo = *incompleteTodo
		return nil
	}

	TodoRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
		updatedTodo = todo // Capture the update
		return nil
	}

	var capturedEvent *model.TodoCompletedEvent
	TodoProducer.SendFunc = func(ctx context.Context, event *model.TodoCompletedEvent) error {
		capturedEvent = event
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Complete(context.Background(), req)

	// ------------------------------------------------------- //

	// Verify todo was updated correctly
	assert.NotNil(t, updatedTodo)
	assert.True(t, updatedTodo.IsCompleted)
	assert.NotNil(t, updatedTodo.CompletedAt)
	assert.Equal(t, incompleteTodo.ID, updatedTodo.ID)
	assert.Equal(t, incompleteTodo.UserID, updatedTodo.UserID)
	assert.Equal(t, incompleteTodo.Title, updatedTodo.Title)
	assert.Equal(t, incompleteTodo.Description, updatedTodo.Description)

	// Verify response matches updated todo
	assert.Equal(t, updatedTodo.ID, res.ID)
	assert.Equal(t, updatedTodo.Title, res.Title)
	assert.Equal(t, updatedTodo.Description, res.Description)
	assert.Equal(t, updatedTodo.IsCompleted, res.IsCompleted)
	assert.Equal(t, updatedTodo.CompletedAt, res.CompletedAt)
	assert.Equal(t, updatedTodo.CreatedAt, res.CreatedAt)
	assert.Equal(t, updatedTodo.UpdatedAt, res.UpdatedAt)

	// Verify event was published
	assert.NotNil(t, capturedEvent)
	assert.Equal(t, updatedTodo.ID, capturedEvent.ID)
	assert.Equal(t, updatedTodo.Title, capturedEvent.Title)
	assert.Equal(t, updatedTodo.UserID, capturedEvent.UserID)
	assert.Equal(t, updatedTodo.CompletedAt, capturedEvent.CompletedAt)

	assert.Nil(t, err)
}

func TestTodoUsecaseImpl_Complete_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.CompleteTodoRequest{
		// Missing required fields UserID and ID
	}

	// ------------------------------------------------------- //

	res, err := u.Complete(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Field validation")
}

func TestTodoUsecaseImpl_Complete_Fail_FindByIDAndUserID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.CompleteTodoRequest{
		ID:     11,
		UserID: testUserID,
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id int64, userID int64) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Complete(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}

func TestTodoUsecaseImpl_Complete_AlreadyCompleted(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	const todoID int64 = 12
	req := &model.CompleteTodoRequest{
		ID:     todoID,
		UserID: testUserID,
	}

	completedAt := time.Now().UnixMilli()
	completedTodo := &entity.Todo{
		ID:          todoID,
		UserID:      req.UserID,
		Title:       "Test Todo",
		Description: "Test Description",
		IsCompleted: true,
		CompletedAt: &completedAt,
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id int64, userID int64) error {
		*todo = *completedTodo
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Complete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.True(t, res.IsCompleted)
	assert.Equal(t, completedAt, *res.CompletedAt)
	assert.Equal(t, completedTodo.Title, res.Title)
	assert.Equal(t, completedTodo.Description, res.Description)
	assert.Nil(t, err)
}

func TestTodoUsecaseImpl_Complete_Fail_Update(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	const todoID int64 = 13
	req := &model.CompleteTodoRequest{
		ID:     todoID,
		UserID: testUserID,
	}

	incompleteTodo := &entity.Todo{
		ID:          todoID,
		UserID:      req.UserID,
		Title:       "Test Todo",
		Description: "Test Description",
		IsCompleted: false,
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id int64, userID int64) error {
		*todo = *incompleteTodo
		return nil
	}

	TodoRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Complete(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}

func TestTodoUsecaseImpl_Complete_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	TodoProducer := &mock.TodoProducerMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
		TodoProducer:   TodoProducer,
	}

	// ------------------------------------------------------- //

	const todoID int64 = 14
	req := &model.CompleteTodoRequest{
		ID:     todoID,
		UserID: testUserID,
	}

	incompleteTodo := &entity.Todo{
		ID:          todoID,
		UserID:      req.UserID,
		Title:       "Test Todo",
		Description: "Test Description",
		IsCompleted: false,
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id int64, userID int64) error {
		*todo = *incompleteTodo
		return nil
	}

	TodoRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
		return nil
	}

	TodoProducer.SendFunc = func(ctx context.Context, event *model.TodoCompletedEvent) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Complete(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}
