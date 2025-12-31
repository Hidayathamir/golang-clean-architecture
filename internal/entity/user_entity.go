package entity

import "time"

type User struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Username  string    `gorm:"column:username;uniqueIndex"`
	Password  string    `gorm:"column:password"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`

	Contacts ContactList `gorm:"foreignKey:user_id;references:id"`
	Todos    TodoList    `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}

type UserList []User
