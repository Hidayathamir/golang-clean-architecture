package todo_test

import (
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewTodoUsecase(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	var DB = gormDB
	var Config = viper.New()

	var TodoRepository repository.TodoRepository = &mock.TodoRepositoryMock{}
	var TodoProducer messaging.TodoProducer = &mock.TodoProducerMock{}

	u := todo.NewTodoUsecase(
		Config, DB,
		TodoRepository,
		TodoProducer,
	)

	assert.NotEmpty(t, u)
}
