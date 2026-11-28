package chatusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ChatUsecaseImpl) ListConversations(ctx context.Context) (dto.ConversationsResponse, error) {
	userAuth := ctxuserauth.Get(ctx)
	conversations := entity.ConversationList{}
	err := u.ChatRepository.ListConversations(ctx, u.DB, &conversations, userAuth.ID)
	if err != nil {
		return dto.ConversationsResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).ListConversations")
	}
	seen := map[int64]struct{}{}
	ids := make([]int64, 0, len(conversations))
	for i := range conversations {
		peerID := conversations[i].PeerID
		if _, ok := seen[peerID]; ok {
			continue
		}
		seen[peerID] = struct{}{}
		ids = append(ids, peerID)
	}
	nameByID := map[int64]string{}
	usernameByID := map[int64]string{}
	if len(ids) > 0 {
		userList := entity.UserList{}
		err = u.UserRepository.FindByIDs(ctx, u.DB, &userList, ids)
		if err != nil {
			return dto.ConversationsResponse{}, errkit.AddFuncName(err, "chatusecase.(*ChatUsecaseImpl).ListConversations")
		}
		for i := range userList {
			nameByID[userList[i].ID] = userList[i].Name
			usernameByID[userList[i].ID] = userList[i].Username
		}
	}
	res := dto.ConversationResponseList{}
	for i := range conversations {
		conversation := conversations[i]
		lastChat := dto.ChatResponse{}
		converter.EntityChatToDtoChatResponse(conversation.LastChat, &lastChat)
		res = append(res, dto.ConversationResponse{
			PeerID:       conversation.PeerID,
			PeerName:     nameByID[conversation.PeerID],
			PeerUsername: usernameByID[conversation.PeerID],
			LastChat:     lastChat,
			UnreadCount:  conversation.UnreadCount,
		})
	}
	return dto.ConversationsResponse{Data: res}, nil
}
