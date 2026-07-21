package topic

const (
	ImageUploaded  = "image.uploaded"
	ImageLiked     = "image.liked"
	ImageCommented = "image.commented"
	UserFollowed   = "user.followed"
	Notif          = "notif"

	ImageUploadedRetry  = "image.uploaded.retry"
	ImageLikedRetry     = "image.liked.retry"
	ImageCommentedRetry = "image.commented.retry"
	UserFollowedRetry   = "user.followed.retry"
	NotifRetry          = "notif.retry"

	ImageUploadedDLQ  = "image.uploaded.dlq"
	ImageLikedDLQ     = "image.liked.dlq"
	ImageCommentedDLQ = "image.commented.dlq"
	UserFollowedDLQ   = "user.followed.dlq"
	NotifDLQ          = "notif.dlq"
)

var PrimaryToRetry = map[string]string{
	ImageUploaded:  ImageUploadedRetry,
	ImageLiked:     ImageLikedRetry,
	ImageCommented: ImageCommentedRetry,
	UserFollowed:   UserFollowedRetry,
	Notif:          NotifRetry,
}

var PrimaryToDLQ = map[string]string{
	ImageUploaded:  ImageUploadedDLQ,
	ImageLiked:     ImageLikedDLQ,
	ImageCommented: ImageCommentedDLQ,
	UserFollowed:   UserFollowedDLQ,
	Notif:          NotifDLQ,
}
