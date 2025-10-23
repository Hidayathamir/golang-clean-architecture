package entity

// Todo represents a task owned by a user.
type Todo struct {
	ID          string `gorm:"column:id;primaryKey"`
	UserID      string `gorm:"column:user_id"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	IsCompleted bool   `gorm:"column:is_completed"`
	CompletedAt *int64 `gorm:"column:completed_at"`
	CreatedAt   int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`

	User User `gorm:"foreignKey:user_id;references:id"`
}

func (t *Todo) TableName() string {
	return "todos"
}
