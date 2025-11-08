package todo_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/stretchr/testify/assert"
)

func TestNewTodoUsecaseMwLogger(t *testing.T) {
	next := &mock.TodoUsecaseMock{}

	mw := todo.NewTodoUsecaseMwLogger(next)

	assert.NotEmpty(t, mw)
}

func TestTodoUsecaseMwLogger_Create(t *testing.T) {
	next := &mock.TodoUsecaseMock{}
	u := todo.NewTodoUsecaseMwLogger(next)

	req := &model.CreateTodoRequest{}
	res := &model.TodoResponse{}
	next.CreateFunc = func(ctx context.Context, req *model.CreateTodoRequest) (*model.TodoResponse, error) {
		return res, nil
	}

	actualRes, err := u.Create(context.Background(), req)

	assert.Equal(t, res, actualRes)
	assert.Nil(t, err)
}

func TestTodoUsecaseMwLogger_Get(t *testing.T) {
	next := &mock.TodoUsecaseMock{}
	u := todo.NewTodoUsecaseMwLogger(next)

	req := &model.GetTodoRequest{}
	res := &model.TodoResponse{}
	next.GetFunc = func(ctx context.Context, req *model.GetTodoRequest) (*model.TodoResponse, error) {
		return res, nil
	}

	actualRes, err := u.Get(context.Background(), req)

	assert.Equal(t, res, actualRes)
	assert.Nil(t, err)
}

func TestTodoUsecaseMwLogger_List(t *testing.T) {
	next := &mock.TodoUsecaseMock{}
	u := todo.NewTodoUsecaseMwLogger(next)

	req := &model.ListTodoRequest{}
	res := model.TodoResponseList{}
	var total int64 = 0
	next.ListFunc = func(ctx context.Context, req *model.ListTodoRequest) (model.TodoResponseList, int64, error) {
		return res, total, nil
	}

	actualRes, actualTotal, err := u.List(context.Background(), req)

	assert.Equal(t, res, actualRes)
	assert.Equal(t, total, actualTotal)
	assert.Nil(t, err)
}

func TestTodoUsecaseMwLogger_Update(t *testing.T) {
	next := &mock.TodoUsecaseMock{}
	u := todo.NewTodoUsecaseMwLogger(next)

	req := &model.UpdateTodoRequest{}
	res := &model.TodoResponse{}
	next.UpdateFunc = func(ctx context.Context, req *model.UpdateTodoRequest) (*model.TodoResponse, error) {
		return res, nil
	}

	actualRes, err := u.Update(context.Background(), req)

	assert.Equal(t, res, actualRes)
	assert.Nil(t, err)
}

func TestTodoUsecaseMwLogger_Delete(t *testing.T) {
	next := &mock.TodoUsecaseMock{}
	u := todo.NewTodoUsecaseMwLogger(next)

	req := &model.DeleteTodoRequest{}
	next.DeleteFunc = func(ctx context.Context, req *model.DeleteTodoRequest) error {
		return nil
	}

	err := u.Delete(context.Background(), req)

	assert.Nil(t, err)
}

func TestTodoUsecaseMwLogger_Complete(t *testing.T) {
	next := &mock.TodoUsecaseMock{}
	u := todo.NewTodoUsecaseMwLogger(next)

	req := &model.CompleteTodoRequest{}
	res := &model.TodoResponse{}
	next.CompleteFunc = func(ctx context.Context, req *model.CompleteTodoRequest) (*model.TodoResponse, error) {
		return res, nil
	}

	actualRes, err := u.Complete(context.Background(), req)

	assert.Equal(t, res, actualRes)
	assert.Nil(t, err)
}
