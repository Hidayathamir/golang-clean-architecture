package entity

type User2 struct {
	ID          string `gorm:"column:id;primaryKey"`
	Email       string `gorm:"column:email"`
	Password    string `gorm:"column:password_hash"`
	DisplayName string `gorm:"column:display_name"`
	CreatedAt   int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (User2) TableName() string {
	return "user2_users"
}

type User2List []User2
