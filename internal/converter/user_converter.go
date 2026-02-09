package converter

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
)

func DtoRegisterUserRequestToEntityUser(req *dto.RegisterUserRequest, user *entity.User, password string) {
	if req == nil || user == nil {
		return
	}

	user.Username = req.Username
	user.Name = req.Name
	user.Password = password
}

func DtoUpdateUserRequestToEntityUser(req *dto.UpdateUserRequest, user *entity.User, password string) {
	if req == nil || user == nil {
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if password != "" {
		user.Password = password
	}
}

func EntityUserToDtoUserResponse(user *entity.User, res *dto.UserResponse) {
	res.ID = user.ID
	res.Username = user.Username
	res.Name = user.Name
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
}

func EntityUserToDtoUserLoginResponse(user *entity.User, res *dto.UserLoginResponse) {
	res.ID = user.ID
	res.Username = user.Username
	res.Name = user.Name
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
}

func EntityUserToDtoUserAuth(user *entity.User, userAuth *dto.UserAuth) {
	if user == nil || userAuth == nil {
		return
	}

	userAuth.ID = user.ID
	userAuth.Username = user.Username
	userAuth.Name = user.Name
	userAuth.FollowerCount = user.FollowerCount
	userAuth.FollowingCount = user.FollowingCount
	userAuth.CreatedAt = user.CreatedAt
	userAuth.UpdatedAt = user.UpdatedAt
	userAuth.DeletedAt = user.DeletedAt
}

func DtoFollowUserRequestToEntityFollow(ctx context.Context, req *dto.FollowUserRequest, follow *entity.Follow) {
	userAuth := ctxuserauth.Get(ctx)
	follow.FollowerID = userAuth.ID
	follow.FollowingID = req.FollowingID
}

func EntityFollowToDtoUserFollowedEvent(ctx context.Context, follow *entity.Follow, event *dto.UserFollowedEvent) {
	event.ID = follow.ID
	event.FollowerID = follow.FollowerID
	event.FollowingID = follow.FollowingID
	event.CreatedAt = follow.CreatedAt
	event.UpdatedAt = follow.UpdatedAt
	event.DeletedAt = follow.DeletedAt
}

func DtoUserFollowedEventToDtoNotifyUserBeingFollowedRequest(ctx context.Context, event *dto.UserFollowedEvent, req *dto.NotifyUserBeingFollowedRequest) {
	req.FollowerID = event.FollowerID
	req.FollowingID = event.FollowingID
}

func SaramaConsumerMessageListToDtoBatchUpdateUserFollowStatsRequest(ctx context.Context, messages []*sarama.ConsumerMessage, req *dto.BatchUpdateUserFollowStatsRequest) {
	userFollowerCounts := make(map[int64]int)
	userFollowingCounts := make(map[int64]int)

	for _, message := range messages {
		event := new(dto.UserFollowedEvent)
		if err := json.Unmarshal(message.Value, event); err != nil {
			x.Logger.WithContext(ctx).WithError(err).Warn("Failed to unmarshal user followed event")
			continue
		}
		userFollowerCounts[event.FollowingID]++
		userFollowingCounts[event.FollowerID]++
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
