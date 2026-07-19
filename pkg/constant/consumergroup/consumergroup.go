// Consumer group naming convention:
// <topic>.<role>-<detail>
//
// Roles:
//
//	notify - per-record, real-time notification delivery
//	batch  - aggregated batch counter/stat updates
//	sync   - one-way data sync to external systems
//	log    - debugging/dummy consumer
package consumergroup

const (
	ImageUploadedNotifyFollowers = "image.uploaded.notify-followers"
	ImageUploadedSyncSearch      = "image.uploaded.sync-search"
	ImageLikedNotifyOwner        = "image.liked.notify-owner"
	ImageLikedBatchCount         = "image.liked.batch-count"
	ImageCommentedNotifyOwner    = "image.commented.notify-owner"
	ImageCommentedBatchCount     = "image.commented.batch-count"

	UserFollowedNotifyUser = "user.followed.notify-user"
	UserFollowedBatchStats = "user.followed.batch-stats"

	NotifLog = "notif.log"
)
