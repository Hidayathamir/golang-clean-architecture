package entity

import (
	"time"

	"github.com/guregu/null/v6"
)

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
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID      int64     `gorm:"column:user_id"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`
	IsCompleted bool      `gorm:"column:is_completed"`
	CompletedAt null.Time `gorm:"column:completed_at;type:timestamptz"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`

	User User `gorm:"foreignKey:user_id;references:id"`
}

func (t *Todo) TableName() string {
	return TodoTableName
}

type TodoList []Todo
