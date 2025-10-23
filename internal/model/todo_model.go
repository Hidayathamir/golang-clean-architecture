package model

type TodoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	CompletedAt *int64 `json:"completed_at,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type CreateTodoRequest struct {
	UserID      string `json:"-"           validate:"required"`
	Title       string `json:"title"       validate:"required,max=200"`
	Description string `json:"description" validate:"max=1000"`
}

type UpdateTodoRequest struct {
	UserID      string `json:"-"           validate:"required"`
	ID          string `json:"-"           validate:"required,max=100,uuid"`
	Title       string `json:"title"       validate:"required,max=200"`
	Description string `json:"description" validate:"max=1000"`
}

type GetTodoRequest struct {
	UserID string `json:"-" validate:"required"`
	ID     string `json:"-" validate:"required,max=100,uuid"`
}

type DeleteTodoRequest struct {
	UserID string `json:"-" validate:"required"`
	ID     string `json:"-" validate:"required,max=100,uuid"`
}

type CompleteTodoRequest struct {
	UserID string `json:"-" validate:"required"`
	ID     string `json:"-" validate:"required,max=100,uuid"`
}

type ListTodoRequest struct {
	UserID      string `json:"-"            validate:"required"`
	Title       string `json:"title"        validate:"max=200"`
	IsCompleted *bool  `json:"is_completed"`
	Page        int    `json:"page"         validate:"min=1"`
	Size        int    `json:"size"         validate:"min=1,max=100"`
}
