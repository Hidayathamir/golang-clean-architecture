package dto

import "time"

type SendChatRequest struct {
	RecipientID int64  `json:"recipient_id" validate:"required"`
	Content     string `json:"content"      validate:"required,max=2000"`
}

type ChatResponse struct {
	ID          int64      `json:"id"`
	SenderID    int64      `json:"sender_id"`
	RecipientID int64      `json:"recipient_id"`
	Content     string     `json:"content"`
	ReadAt      *time.Time `json:"read_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type ChatResponseList []ChatResponse

type GetThreadRequest struct {
	PeerID int64 `json:"peer_id" validate:"required"`
	Page   PageRequest
}

type ThreadResponse = CursorPage[ChatResponseList]

type ConversationResponse struct {
	PeerID       int64        `json:"peer_id"`
	PeerName     string       `json:"peer_name"`
	PeerUsername string       `json:"peer_username"`
	LastChat     ChatResponse `json:"last_chat"`
	UnreadCount  int64        `json:"unread_count"`
}

type ConversationResponseList []ConversationResponse

type ConversationsResponse struct {
	Data ConversationResponseList `json:"data"`
}

type MarkChatReadRequest struct {
	PeerID int64 `json:"peer_id" validate:"required"`
}

type UnreadChatCountResponse struct {
	UnreadCount int64 `json:"unread_count"`
}
