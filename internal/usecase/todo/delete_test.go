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

func TestTodoUsecaseImpl_Delete_Success(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectCommit()

	req := &model.DeleteTodoRequest{
		ID:     uuid.NewString(),
		UserID: "user1",
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id, userID string) error {
		return nil
	}

	TodoRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
		return nil
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.Nil(t, err)
}

func TestTodoUsecaseImpl_Delete_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.DeleteTodoRequest{
		// Missing required fields
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
}

func TestTodoUsecaseImpl_Delete_Fail_FindByIDAndUserID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	req := &model.DeleteTodoRequest{
		ID:     uuid.NewString(),
		UserID: "user1",
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id, userID string) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
}

func TestTodoUsecaseImpl_Delete_Fail_Delete(t *testing.T) {
	gormDB, sqlMockDB := newFakeDB(t)
	TodoRepository := &mock.TodoRepositoryMock{}
	u := &todo.TodoUsecaseImpl{
		DB:             gormDB,
		TodoRepository: TodoRepository,
	}

	// ------------------------------------------------------- //

	sqlMockDB.ExpectBegin()
	sqlMockDB.ExpectRollback()

	req := &model.DeleteTodoRequest{
		ID:     uuid.NewString(),
		UserID: "user1",
	}

	TodoRepository.FindByIDAndUserIDFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo, id, userID string) error {
		return nil
	}

	TodoRepository.DeleteFunc = func(ctx context.Context, db *gorm.DB, todo *entity.Todo) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	err := u.Delete(context.Background(), req)

	// ------------------------------------------------------- //

	assert.NotNil(t, err)
}
