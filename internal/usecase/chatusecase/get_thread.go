package chatusecase

import (
	"context"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ChatUsecaseImpl) GetThread(ctx context.Context, req dto.GetThreadRequest) (dto.ThreadResponse, error) {
	userAuth := ctxuserauth.Get(ctx)
	cursor, hasCursor, err := req.Page.DecodeCursor()
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.ThreadResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).GetThread")
	}
	var cursorCreatedAt time.Time
	var cursorID int64
	if hasCursor {
		cursorCreatedAt = cursor.CreatedAt
		cursorID = cursor.ID
	}
	limit := req.Page.NormalizedLimit()
	chatList := entity.ChatList{}
	err = u.ChatRepository.FindThread(ctx, u.DB, &chatList, userAuth.ID, req.PeerID, cursorCreatedAt, cursorID, limit+1)
	if err != nil {
		return dto.ThreadResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).GetThread")
	}
	hasMore := len(chatList) > limit
	if hasMore {
		chatList = chatList[:limit]
	}
	res := dto.ChatResponseList{}
	converter.EntityChatListToDtoChatResponseList(chatList, &res)
	var lastCreatedAt time.Time
	var lastID int64
	if len(chatList) > 0 {
		last := chatList[len(chatList)-1]
		lastCreatedAt = last.CreatedAt
		lastID = last.ID
	}
	return dto.PageFrom(res, hasMore, lastCreatedAt, lastID), nil
}
