package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	todousecase "github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type TodoCommandConsumer struct {
	Usecase todousecase.TodoUsecase
	Log     *logrus.Logger
}

func NewTodoCommandConsumer(usecase todousecase.TodoUsecase, log *logrus.Logger) *TodoCommandConsumer {
	return &TodoCommandConsumer{
		Usecase: usecase,
		Log:     log,
	}
}

type completeTodoCommand struct {
	UserID string `json:"user_id"`
	TodoID string `json:"todo_id"`
}

func (c *TodoCommandConsumer) Consume(message *sarama.ConsumerMessage) error {
	cmd := new(completeTodoCommand)
	if err := json.Unmarshal(message.Value, cmd); err != nil {
		c.Log.WithError(err).Error("error unmarshalling todo command")
		return errkit.AddFuncName("messaging.(*TodoCommandConsumer).Consume", err)
	}

	if cmd.UserID == "" || cmd.TodoID == "" {
		err := errkit.BadRequest(fmt.Errorf("invalid todo command payload: %+v", cmd))
		c.Log.WithError(err).Warn("todo command missing identifiers")
		return errkit.AddFuncName("messaging.(*TodoCommandConsumer).Consume", err)
	}

	req := &model.CompleteTodoRequest{
		UserID: cmd.UserID,
		ID:     cmd.TodoID,
	}

	if _, err := c.Usecase.Complete(context.Background(), req); err != nil {
		return errkit.AddFuncName("messaging.(*TodoCommandConsumer).Consume", err)
	}

	c.Log.WithFields(logrus.Fields{
		"user_id": cmd.UserID,
		"todo_id": cmd.TodoID,
	}).Info("Processed todo completion command")

	return nil
}
