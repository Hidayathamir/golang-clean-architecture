package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	todousecase "github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type TodoCommandConsumer struct {
	Usecase todousecase.TodoUsecase
}

func NewTodoCommandConsumer(usecase todousecase.TodoUsecase) *TodoCommandConsumer {
	return &TodoCommandConsumer{
		Usecase: usecase,
	}
}

type completeTodoCommand struct {
	UserID int64 `json:"user_id"`
	TodoID int64 `json:"todo_id"`
}

func (c *TodoCommandConsumer) Consume(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	cmd := new(completeTodoCommand)
	if err := json.Unmarshal(message.Value, cmd); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error("error unmarshalling todo command")
		return errkit.AddFuncName(err)
	}

	if cmd.UserID <= 0 || cmd.TodoID <= 0 {
		err := errkit.BadRequest(fmt.Errorf("invalid todo command payload: %+v", cmd))
		x.Logger.WithContext(ctx).WithError(err).Warn("todo command missing identifiers")
		return errkit.AddFuncName(err)
	}

	req := &model.CompleteTodoRequest{
		UserID: cmd.UserID,
		ID:     cmd.TodoID,
	}

	if _, err := c.Usecase.Complete(ctx, req); err != nil {
		return errkit.AddFuncName(err)
	}

	x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"user_id": cmd.UserID,
		"todo_id": cmd.TodoID,
	}).Info("Processed todo completion command")

	return nil
}
