package model

import "strconv"

type TodoCompletedEvent struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CompletedAt *int64 `json:"completed_at"`
}

func (t *TodoCompletedEvent) GetID() string {
	return strconv.FormatInt(t.ID, 10)
}
