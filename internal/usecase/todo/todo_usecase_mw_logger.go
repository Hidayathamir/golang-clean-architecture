package todo

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
	"github.com/sirupsen/logrus"
)

var _ TodoUsecase = &TodoUsecaseMwLogger{}

type TodoUsecaseMwLogger struct {
	Next TodoUsecase
}

func NewTodoUsecaseMwLogger(next TodoUsecase) *TodoUsecaseMwLogger {
	return &TodoUsecaseMwLogger{
		Next: next,
	}
}

func (u *TodoUsecaseMwLogger) Create(ctx context.Context, req *model.CreateTodoRequest) (*model.TodoResponse, error) {
	res, err := u.Next.Create(ctx, req)
	logging.Log(ctx, logrus.Fields{"req": req, "res": res}, err)
	return res, err
}

func (u *TodoUsecaseMwLogger) Get(ctx context.Context, req *model.GetTodoRequest) (*model.TodoResponse, error) {
	res, err := u.Next.Get(ctx, req)
	logging.Log(ctx, logrus.Fields{"req": req, "res": res}, err)
	return res, err
}

func (u *TodoUsecaseMwLogger) List(ctx context.Context, req *model.ListTodoRequest) (model.TodoResponseList, int64, error) {
	res, total, err := u.Next.List(ctx, req)
	logging.Log(ctx, logrus.Fields{"req": req, "res": res, "total": total}, err)
	return res, total, err
}

func (u *TodoUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateTodoRequest) (*model.TodoResponse, error) {
	res, err := u.Next.Update(ctx, req)
	logging.Log(ctx, logrus.Fields{"req": req, "res": res}, err)
	return res, err
}

func (u *TodoUsecaseMwLogger) Delete(ctx context.Context, req *model.DeleteTodoRequest) error {
	err := u.Next.Delete(ctx, req)
	logging.Log(ctx, logrus.Fields{"req": req}, err)
	return err
}

func (u *TodoUsecaseMwLogger) Complete(ctx context.Context, req *model.CompleteTodoRequest) (*model.TodoResponse, error) {
	res, err := u.Next.Complete(ctx, req)
	logging.Log(ctx, logrus.Fields{"req": req, "res": res}, err)
	return res, err
}
