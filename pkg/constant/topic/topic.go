package topic

type Topic struct {
	Primary string
}

func (t Topic) Retry() string {
	return t.Primary + ".retry"
}

func (t Topic) DLQ() string {
	return t.Primary + ".dlq"
}

var (
	ImageUploaded  = Topic{Primary: "image.uploaded"}
	ImageLiked     = Topic{Primary: "image.liked"}
	ImageCommented = Topic{Primary: "image.commented"}
	UserFollowed   = Topic{Primary: "user.followed"}
	Notif          = Topic{Primary: "notif"}
)
