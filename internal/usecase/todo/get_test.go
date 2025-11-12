package todo_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTodoUsecaseImpl_Get_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetTodoRequest{
		ID:     uuid.NewString(),
		UserID: testUserID,
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id string, userID int64) error {
		todo.ID = req.ID
		todo.UserID = req.UserID
		todo.Title = "Test Todo"
		todo.Description = "Test Description"
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Get(context.Background(), req)

	// ------------------------------------------------------- //

	expected := &model.TodoResponse{
		ID:          req.ID,
		Title:       "Test Todo",
		Description: "Test Description",
	}

	assert.Equal(t, expected.ID, res.ID)
	assert.Equal(t, expected.Title, res.Title)
	assert.Equal(t, expected.Description, res.Description)
	assert.Nil(t, err)
}

func TestTodoUsecaseImpl_Get_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetTodoRequest{
		// Missing required fields
	}

	// ------------------------------------------------------- //

	res, err := u.Get(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}

func TestTodoUsecaseImpl_Get_Fail_FindByIDAndUserID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.GetTodoRequest{
		ID:     uuid.NewString(),
		UserID: testUserID,
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id string, userID int64) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Get(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}
