package entity

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"gorm.io/gorm"
)

type Chat struct {
	ID          int64          `gorm:"column:id;primaryKey"`
	SenderID    int64          `gorm:"column:sender_id"`
	RecipientID int64          `gorm:"column:recipient_id"`
	Content     string         `gorm:"column:content"`
	ReadAt      *time.Time     `gorm:"column:read_at"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (c *Chat) TableName() string {
	return table.Chat
}

func (c *Chat) IsValid() bool {
	return c.Content != ""
}

type ChatList []Chat

// Conversation is one thread's summary from the current user's perspective:
// the peer, the latest chat, and how many of the peer's chats are unread.
type Conversation struct {
	PeerID       int64
	PeerName     string
	PeerUsername string
	LastChat     Chat
	LastChatAt   time.Time
	UnreadCount  int64
}

type ConversationList []Conversation
