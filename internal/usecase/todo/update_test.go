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

func TestTodoUsecaseImpl_Update_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	const title = "Updated Title"
	req := &model.UpdateTodoRequest{
		ID:          uuid.NewString(),
		UserID:      "user1",
		Title:       title,
		Description: "Updated Description",
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id, userID string) error {
		todo.ID = req.ID
		todo.UserID = req.UserID
		return nil
	}

	TodoRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	expected := &model.TodoResponse{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
	}

	assert.Equal(t, expected.ID, res.ID)
	assert.Equal(t, expected.Title, res.Title)
	assert.Equal(t, expected.Description, res.Description)
	assert.Nil(t, err)
}

func TestTodoUsecaseImpl_Update_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.UpdateTodoRequest{
		// Missing required fields
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}

func TestTodoUsecaseImpl_Update_Fail_FindByIDAndUserID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.UpdateTodoRequest{
		ID:          uuid.NewString(),
		UserID:      "user1",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id, userID string) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}

func TestTodoUsecaseImpl_Update_Fail_Update(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.UpdateTodoRequest{
		ID:          uuid.NewString(),
		UserID:      "user1",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id, userID string) error {
		return nil
	}

	TodoRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.TodoResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}
