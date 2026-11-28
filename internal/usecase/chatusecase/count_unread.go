package chatusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ChatUsecaseImpl) CountUnread(ctx context.Context) (dto.UnreadChatCountResponse, error) {
	userAuth := ctxuserauth.Get(ctx)
	total, err := u.ChatRepository.CountUnread(ctx, u.DB, userAuth.ID)
	if err != nil {
		return dto.UnreadChatCountResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).CountUnread")
	}
	return dto.UnreadChatCountResponse{UnreadCount: total}, nil
}
