package model

type TodoResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	CompletedAt *int64 `json:"completed_at,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type CreateTodoRequest struct {
	UserID      int64  `json:"-"           validate:"required"`
	Title       string `json:"title"       validate:"required,max=200"`
	Description string `json:"description" validate:"max=1000"`
}

type UpdateTodoRequest struct {
	UserID      int64  `json:"-"           validate:"required"`
	ID          int64  `json:"-"           validate:"required,min=1"`
	Title       string `json:"title"       validate:"required,max=200"`
	Description string `json:"description" validate:"max=1000"`
}

type GetTodoRequest struct {
	UserID int64 `json:"-" validate:"required"`
	ID     int64 `json:"-" validate:"required,min=1"`
}

type DeleteTodoRequest struct {
	UserID int64 `json:"-" validate:"required"`
	ID     int64 `json:"-" validate:"required,min=1"`
}

type CompleteTodoRequest struct {
	UserID int64 `json:"-" validate:"required"`
	ID     int64 `json:"-" validate:"required,min=1"`
}

type ListTodoRequest struct {
	UserID      int64  `json:"-"            validate:"required"`
	Title       string `json:"title"        validate:"max=200"`
	IsCompleted *bool  `json:"is_completed"`
	Page        int    `json:"page"         validate:"min=1"`
	Size        int    `json:"size"         validate:"min=1,max=100"`
}

type TodoResponseList []TodoResponse
