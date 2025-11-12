package entity

const (
	TodoTableName         = "todos"
	TodoColumnID          = "id"
	TodoColumnUserID      = "user_id"
	TodoColumnTitle       = "title"
	TodoColumnDescription = "description"
	TodoColumnIsCompleted = "is_completed"
	TodoColumnCompletedAt = "completed_at"
	TodoColumnCreatedAt   = "created_at"
	TodoColumnUpdatedAt   = "updated_at"
)

type Todo struct {
	ID          string `gorm:"column:id;primaryKey"` // TODO: this is uuid, next use int
	UserID      int64  `gorm:"column:user_id"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	IsCompleted bool   `gorm:"column:is_completed"`
	CompletedAt *int64 `gorm:"column:completed_at"`
	CreatedAt   int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`

	User User `gorm:"foreignKey:user_id;references:id"`
}

func (t *Todo) TableName() string {
	return TodoTableName
}

type TodoList []Todo
