package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
)

func EntityChatToDtoChatResponse(chat entity.Chat, res *dto.ChatResponse) {
	res.ID = chat.ID
	res.SenderID = chat.SenderID
	res.RecipientID = chat.RecipientID
	res.Content = chat.Content
	res.ReadAt = chat.ReadAt
	res.CreatedAt = chat.CreatedAt
}

func EntityChatListToDtoChatResponseList(chatList entity.ChatList, res *dto.ChatResponseList) {
	for i := range chatList {
		item := dto.ChatResponse{}
		EntityChatToDtoChatResponse(chatList[i], &item)
		*res = append(*res, item)
	}
}
