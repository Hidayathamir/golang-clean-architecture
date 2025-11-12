package entity

type User struct {
	ID        int64  `gorm:"column:id;primaryKey"`
	Username  string `gorm:"column:username;uniqueIndex"`
	Password  string `gorm:"column:password"`
	Name      string `gorm:"column:name"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`

	Contacts ContactList `gorm:"foreignKey:user_id;references:id"`
	Todos    TodoList    `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}

type UserList []User
