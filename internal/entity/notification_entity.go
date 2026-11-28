package entity

import (
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/table"
	"gorm.io/gorm"
)

const (
	NotificationTypeFollow  = "follow"
	NotificationTypeLike    = "like"
	NotificationTypeComment = "comment"
	NotificationTypeNewPost = "new_post"
)

type Notification struct {
	ID        int64          `gorm:"column:id;primaryKey"`
	UserID    int64          `gorm:"column:user_id"`
	ActorID   int64          `gorm:"column:actor_id"`
	Type      string         `gorm:"column:type"`
	ImageID   int64          `gorm:"column:image_id"`
	IsRead    bool           `gorm:"column:is_read"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (n *Notification) TableName() string {
	return table.Notification
}

func (n *Notification) IsValid() bool {
	switch n.Type {
	case NotificationTypeFollow, NotificationTypeLike, NotificationTypeComment, NotificationTypeNewPost:
		return true
	default:
		return false
	}
}

type NotificationList []Notification
