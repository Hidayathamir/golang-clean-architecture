package chatusecase

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ChatUsecaseImpl) MarkRead(ctx context.Context, req dto.MarkChatReadRequest) error {
	userAuth := ctxuserauth.Get(ctx)
	_, err := u.ChatRepository.MarkThreadRead(ctx, u.DB, userAuth.ID, req.PeerID, time.Now())
	if err != nil {
		return errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).MarkRead")
	}
	return nil
}
