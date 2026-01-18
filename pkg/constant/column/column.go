package column

type Column string

func (c Column) Str() string {
	return string(c)
}

func (c Column) Eq(value any) (string, any) {
	return string(c) + " = ?", value
}

func (c Column) Plus(value any) (string, any) {
	return string(c) + " + ?", value
}

const (
	ID             Column = "id"
	UserID         Column = "user_id"
	ImageID        Column = "image_id"
	Comment        Column = "comment"
	CreatedAt      Column = "created_at"
	UpdatedAt      Column = "updated_at"
	DeletedAt      Column = "deleted_at"
	FollowerID     Column = "follower_id"
	FollowingID    Column = "following_id"
	URL            Column = "url"
	LikeCount      Column = "like_count"
	CommentCount   Column = "comment_count"
	Username       Column = "username"
	Password       Column = "password"
	Name           Column = "name"
	FollowerCount  Column = "follower_count"
	FollowingCount Column = "following_count"
)
