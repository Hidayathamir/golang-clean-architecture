package todo

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/TodoUsecase.go -pkg=mock . TodoUsecase

type TodoUsecase interface {
	Create(ctx context.Context, req *model.CreateTodoRequest) (*model.TodoResponse, error)
	Get(ctx context.Context, req *model.GetTodoRequest) (*model.TodoResponse, error)
	List(ctx context.Context, req *model.ListTodoRequest) (model.TodoResponseList, int64, error)
	Update(ctx context.Context, req *model.UpdateTodoRequest) (*model.TodoResponse, error)
	Delete(ctx context.Context, req *model.DeleteTodoRequest) error
	Complete(ctx context.Context, req *model.CompleteTodoRequest) (*model.TodoResponse, error)
}

var _ TodoUsecase = &TodoUsecaseImpl{}

type TodoUsecaseImpl struct {
	Config   *viper.Viper
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate

	// repository
	TodoRepository repository.TodoRepository

	// producer
	TodoProducer messaging.TodoProducer
}

func NewTodoUsecase(
	cfg *viper.Viper, log *logrus.Logger, db *gorm.DB, validate *validator.Validate,

	// repository
	todoRepository repository.TodoRepository,

	// producer
	todoProducer messaging.TodoProducer,
) *TodoUsecaseImpl {
	return &TodoUsecaseImpl{
		Config:   cfg,
		Log:      log,
		DB:       db,
		Validate: validate,

		// repository
		TodoRepository: todoRepository,

		// producer
		TodoProducer: todoProducer,
	}
}
