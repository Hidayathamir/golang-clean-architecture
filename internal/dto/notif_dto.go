package dto

type NotifyRequest struct {
	UserID  int64  `json:"user_id"`
	ActorID int64  `json:"actor_id"`
	Type    string `json:"type"`
	ImageID int64  `json:"image_id"`
	Message string `json:"message"`
}

type NotifEvent struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}
