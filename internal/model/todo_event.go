package model

type TodoCompletedEvent struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CompletedAt *int64 `json:"completed_at"`
}

func (t *TodoCompletedEvent) GetID() string {
	return t.ID
}
