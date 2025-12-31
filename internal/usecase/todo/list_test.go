package todo_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTodoUsecaseImpl_List_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.ListTodoRequest{
		UserID: testUserID,
		Page:   1,
		Size:   10,
	}

	now := time.UnixMilli(1699432800000).UTC() // Fixed timestamp for test
	mockTodos := entity.TodoList{
		{
			ID:          1,
			UserID:      req.UserID,
			Title:       "Test Todo 1",
			Description: "Description 1",
			IsCompleted: true,
			CompletedAt: null.TimeFrom(now),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			UserID:      req.UserID,
			Title:       "Test Todo 2",
			Description: "Description 2",
			IsCompleted: false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	TodoRepository.ListFunc = func(ctx context.Context, db *gorm.DB, listReq *model.ListTodoRequest) (entity.TodoList, int64, error) {
		assert.Equal(t, req.UserID, listReq.UserID)
		assert.Equal(t, req.Page, listReq.Page)
		assert.Equal(t, req.Size, listReq.Size)
		return mockTodos, int64(len(mockTodos)), nil
	}

	// ------------------------------------------------------- //

	res, total, err := u.List(context.Background(), req)

	// ------------------------------------------------------- //

	assert.Equal(t, 2, len(res))
	assert.Equal(t, int64(2), total)

	// Verify first todo
	assert.Equal(t, mockTodos[0].ID, res[0].ID)
	assert.Equal(t, mockTodos[0].Title, res[0].Title)
	assert.Equal(t, mockTodos[0].Description, res[0].Description)
	assert.Equal(t, mockTodos[0].IsCompleted, res[0].IsCompleted)
	assert.Equal(t, mockTodos[0].CompletedAt, res[0].CompletedAt)
	assert.Equal(t, mockTodos[0].CreatedAt, res[0].CreatedAt)
	assert.Equal(t, mockTodos[0].UpdatedAt, res[0].UpdatedAt)

	// Verify second todo
	assert.Equal(t, mockTodos[1].ID, res[1].ID)
	assert.Equal(t, mockTodos[1].Title, res[1].Title)
	assert.Equal(t, mockTodos[1].Description, res[1].Description)
	assert.Equal(t, mockTodos[1].IsCompleted, res[1].IsCompleted)
	assert.Equal(t, mockTodos[1].CompletedAt, res[1].CompletedAt)
	assert.Equal(t, mockTodos[1].CreatedAt, res[1].CreatedAt)
	assert.Equal(t, mockTodos[1].UpdatedAt, res[1].UpdatedAt)

	assert.Nil(t, err)
}

func TestTodoUsecaseImpl_List_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.ListTodoRequest{
		// Missing required fields UserID, Page, and Size
		Title: "Test",
	}

	// ------------------------------------------------------- //

	res, total, err := u.List(context.Background(), req)

	// ------------------------------------------------------- //

	var expected model.TodoResponseList
	var expectedTotal int64

	assert.Equal(t, expected, res)
	assert.Equal(t, expectedTotal, total)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Field validation")
}

func TestTodoUsecaseImpl_List_Fail_List(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.ListTodoRequest{
		UserID: testUserID,
	}

	TodoRepository.ListFunc = func(ctx context.Context, db *gorm.DB, req *model.ListTodoRequest) (entity.TodoList, int64, error) {
		return nil, 0, assert.AnError
	}

	// ------------------------------------------------------- //

	res, total, err := u.List(context.Background(), req)

	// ------------------------------------------------------- //

	var expected model.TodoResponseList
	var expectedTotal int64

	assert.Equal(t, expected, res)
	assert.Equal(t, expectedTotal, total)
	assert.NotNil(t, err)
}
