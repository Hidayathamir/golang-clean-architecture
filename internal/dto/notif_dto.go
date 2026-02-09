package dto

type NotifyRequest struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

type NotifEvent struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}
