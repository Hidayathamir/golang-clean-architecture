package chatusecase

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
)

func (u *ChatUsecaseImpl) Send(ctx context.Context, req dto.SendChatRequest) (dto.ChatResponse, error) {
	userAuth := ctxuserauth.Get(ctx)
	content := strings.TrimSpace(req.Content)
	if content == "" {
		err := errkit.SetCode(errors.New("content is required"), http.StatusBadRequest)
		return dto.ChatResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).Send")
	}
	if req.RecipientID == userAuth.ID {
		err := errkit.SetCode(errors.New("cannot chat yourself"), http.StatusBadRequest)
		return dto.ChatResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).Send")
	}
	recipient := entity.User{}
	err := u.UserRepository.FindByID(ctx, u.DB, &recipient, req.RecipientID)
	if err != nil {
		return dto.ChatResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).Send")
	}
	chat := entity.Chat{
		SenderID:    userAuth.ID,
		RecipientID: req.RecipientID,
		Content:     content,
	}
	err = u.ChatRepository.Create(ctx, u.DB, &chat)
	if err != nil {
		return dto.ChatResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).Send")
	}
	res := dto.ChatResponse{}
	converter.EntityChatToDtoChatResponse(chat, &res)
	notification := dto.NotificationResponse{
		ID:        chat.ID,
		ActorID:   chat.SenderID,
		ActorName: recipient.Name,
		Type:      "chat",
		CreatedAt: chat.CreatedAt,
	}
	err = u.NotificationPublisher.Publish(ctx, chat.RecipientID, notification)
	if err != nil {
		return dto.ChatResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).Send")
	}
	logkit.Logger.WithContext(ctx).WithField("chat", res).Debug("chat published")
	return res, nil
}
