package column

type Column string

func (c Column) Str() string {
	return string(c)
}

func (c Column) Eq(value any) (string, any) {
	return string(c) + " = ?", value
}

func (c Column) In(value any) (string, any) {
	return string(c) + " IN ?", value
}

// Desc returns the column qualified for a descending ORDER BY.
func (c Column) Desc() string {
	return string(c) + " DESC"
}

// InSubquery returns a "<column> IN (subquery)" predicate for the given subquery.
func (c Column) InSubquery(subquery any) (string, any) {
	return string(c) + " IN (?)", subquery
}

// IsNull returns a "<column> IS NULL" predicate.
func (c Column) IsNull() string {
	return string(c) + " IS NULL"
}

// KeysetBefore returns the keyed-keyset predicate that positions a descending
// (created_at, id) ordered query strictly before the given cursor values.
func KeysetBefore(cursorCreatedAt any, cursorID any) (string, any, any) {
	return "(" + CreatedAt.Str() + ", " + ID.Str() + ") < (?, ?)", cursorCreatedAt, cursorID
}

func (c Column) Plus(value any) (string, any) {
	return string(c) + " + ?", value
}

const (
	ID                Column = "id"
	UserID            Column = "user_id"
	ImageID           Column = "image_id"
	Comment           Column = "comment"
	CreatedAt         Column = "created_at"
	UpdatedAt         Column = "updated_at"
	DeletedAt         Column = "deleted_at"
	FollowerID        Column = "follower_id"
	FollowingID       Column = "following_id"
	URL               Column = "url"
	Caption           Column = "caption"
	LikeCount         Column = "like_count"
	CommentCount      Column = "comment_count"
	Username          Column = "username"
	Password          Column = "password"
	Name              Column = "name"
	FollowerCount     Column = "follower_count"
	FollowingCount    Column = "following_count"
	Topic             Column = "topic"
	Payload           Column = "payload"
	Status            Column = "status"
	ActorID           Column = "actor_id"
	Type              Column = "type"
	IsRead            Column = "is_read"
	RecipientID       Column = "recipient_id"
	SenderID          Column = "sender_id"
	Content           Column = "content"
	ReadAt            Column = "read_at"
	Email             Column = "email"
	EmailVerifiedAt   Column = "email_verified_at"
	VerificationToken Column = "verification_token"
)
