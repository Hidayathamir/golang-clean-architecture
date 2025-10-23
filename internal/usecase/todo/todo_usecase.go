package todo

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/TodoUsecase.go -pkg=mock . TodoUsecase

type TodoUsecase interface {
	Create(ctx context.Context, req *model.CreateTodoRequest) (*model.TodoResponse, error)
	Get(ctx context.Context, req *model.GetTodoRequest) (*model.TodoResponse, error)
	List(ctx context.Context, req *model.ListTodoRequest) ([]model.TodoResponse, int64, error)
	Update(ctx context.Context, req *model.UpdateTodoRequest) (*model.TodoResponse, error)
	Delete(ctx context.Context, req *model.DeleteTodoRequest) error
	Complete(ctx context.Context, req *model.CompleteTodoRequest) (*model.TodoResponse, error)
}

var _ TodoUsecase = &TodoUsecaseImpl{}

type TodoUsecaseImpl struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate

	TodoRepository repository.TodoRepository
	TodoProducer   messaging.TodoProducer
}

func NewTodoUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	todoRepository repository.TodoRepository,
	todoProducer messaging.TodoProducer,
) *TodoUsecaseImpl {
	return &TodoUsecaseImpl{
		DB:             db,
		Log:            log,
		Validate:       validate,
		TodoRepository: todoRepository,
		TodoProducer:   todoProducer,
	}
}
