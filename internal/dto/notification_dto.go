package dto

import "time"

type NotificationResponse struct {
	ID        int64     `json:"id"`
	ActorID   int64     `json:"actor_id"`
	ActorName string    `json:"actor_name"`
	Type      string    `json:"type"`
	ImageID   int64     `json:"image_id"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationResponseList []NotificationResponse

type NotificationsRequest struct {
	Page PageRequest
}

type NotificationsResponse = CursorPage[NotificationResponseList]

type MarkNotificationReadRequest struct {
	NotificationID int64 `json:"notification_id" validate:"required"`
}

type UnreadCountResponse struct {
	UnreadCount int64 `json:"unread_count"`
}
