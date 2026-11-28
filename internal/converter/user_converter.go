package converter

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/twmb/franz-go/pkg/kgo"
)

func DtoRegisterUserRequestToEntityUser(req dto.RegisterUserRequest, user *entity.User, password string) {
	user.Username = req.Username
	user.Name = req.Name
	user.Password = password
}

func DtoUpdateUserRequestToEntityUser(req dto.UpdateUserRequest, user *entity.User, password string) {
	if req.Name != "" {
		user.Name = req.Name
	}

	if password != "" {
		user.Password = password
	}
}

func EntityUserToDtoUserResponse(user entity.User, res *dto.UserResponse) {
	res.ID = user.ID
	res.Username = user.Username
	res.Name = user.Name
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
}

func EntityUserToDtoUserLoginResponse(user entity.User, res *dto.UserLoginResponse) {
	res.ID = user.ID
	res.Username = user.Username
	res.Name = user.Name
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
}

func EntityUserAndUserStatToDtoProfileResponse(user entity.User, userStat entity.UserStat, res *dto.ProfileResponse) {
	res.ID = user.ID
	res.Username = user.Username
	res.Name = user.Name
	res.FollowerCount = userStat.FollowerCount
	res.FollowingCount = userStat.FollowingCount
}

// EntityUserListToDtoUserResponseList maps users in the given order.
func EntityUserListToDtoUserResponseList(userList entity.UserList, res *dto.UserResponseList) {
	*res = make(dto.UserResponseList, 0, len(userList))
	for _, user := range userList {
		item := dto.UserResponse{}
		EntityUserToDtoUserResponse(user, &item)
		*res = append(*res, item)
	}
}

func EntityUserToDtoUserAuth(user entity.User, userAuth *dto.UserAuth) {
	userAuth.ID = user.ID
	userAuth.Username = user.Username
	userAuth.Name = user.Name
	userAuth.CreatedAt = user.CreatedAt
	userAuth.UpdatedAt = user.UpdatedAt
	userAuth.DeletedAt = user.DeletedAt
}

func DtoFollowUserRequestToEntityFollow(ctx context.Context, req dto.FollowUserRequest, follow *entity.Follow) {
	userAuth := ctxuserauth.Get(ctx)
	follow.FollowerID = userAuth.ID
	follow.FollowingID = req.FollowingID
}

func EntityFollowToDtoUserFollowedEvent(follow entity.Follow, event *dto.UserFollowedEvent) {
	event.ID = follow.ID
	event.FollowerID = follow.FollowerID
	event.FollowingID = follow.FollowingID
	event.Count = 1
	event.CreatedAt = follow.CreatedAt
	event.UpdatedAt = follow.UpdatedAt
	event.DeletedAt = follow.DeletedAt
}

func DtoUserFollowedEventToDtoNotifyUserBeingFollowedRequest(event dto.UserFollowedEvent, req *dto.NotifyUserBeingFollowedRequest) {
	req.FollowerID = event.FollowerID
	req.FollowingID = event.FollowingID
}

func KGoRecordListToDtoBatchUpdateUserFollowStatsRequest(ctx context.Context, records []*kgo.Record, req *dto.BatchUpdateUserFollowStatsRequest) {
	userFollowerCounts := make(map[int64]int)
	userFollowingCounts := make(map[int64]int)

	for _, record := range records {
		event := dto.UserFollowedEvent{}
		err := json.Unmarshal(record.Value, &event)
		if err != nil {
			logkit.Logger.WithContext(ctx).WithError(err).Warn("Failed to unmarshal user followed event")
			continue
		}
		// Count is +1 for a follow and -1 for an unfollow. Older events predate
		// the field, so a zero is treated as a single follow.
		count := event.Count
		if count == 0 {
			count = 1
		}
		userFollowerCounts[event.FollowingID] += count
		userFollowingCounts[event.FollowerID] += count
	}

	allUserIDs := make(map[int64]struct{})
	for id := range userFollowerCounts {
		allUserIDs[id] = struct{}{}
	}
	for id := range userFollowingCounts {
		allUserIDs[id] = struct{}{}
	}

	for id := range allUserIDs {
		object := dto.UserIncreaseFollowerFollowingCount{
			UserID:         id,
			FollowerCount:  userFollowerCounts[id],
			FollowingCount: userFollowingCounts[id],
		}
		req.UserIncreaseFollowerFollowingCountList = append(req.UserIncreaseFollowerFollowingCountList, object)
	}
}
