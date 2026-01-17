package converter

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
)

func ModelRegisterUserRequestToEntityUser(req *model.RegisterUserRequest, user *entity.User, password string) {
	if req == nil || user == nil {
		return
	}

	user.Username = req.Username
	user.Name = req.Name
	user.Password = password
}

func ModelUpdateUserRequestToEntityUser(req *model.UpdateUserRequest, user *entity.User, password string) {
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

func EntityUserToModelUserResponse(user *entity.User, res *model.UserResponse) {
	res.ID = user.ID
	res.Username = user.Username
	res.Name = user.Name
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
}

func EntityUserToModelUserLoginResponse(user *entity.User, res *model.UserLoginResponse) {
	res.ID = user.ID
	res.Username = user.Username
	res.Name = user.Name
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
}

func EntityUserToModelUserAuth(user *entity.User, userAuth *model.UserAuth) {
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

func ModelFollowUserRequestToEntityFollow(ctx context.Context, req *model.FollowUserRequest, follow *entity.Follow) {
	userAuth := ctxuserauth.Get(ctx)
	follow.FollowerID = userAuth.ID
	follow.FollowingID = req.FollowingID
}

func EntityFollowToModelUserFollowedEvent(ctx context.Context, follow *entity.Follow, event *model.UserFollowedEvent) {
	event.ID = follow.ID
	event.FollowerID = follow.FollowerID
	event.FollowingID = follow.FollowingID
	event.CreatedAt = follow.CreatedAt
	event.UpdatedAt = follow.UpdatedAt
	event.DeletedAt = follow.DeletedAt
}

func ModelUserFollowedEventToModelNotifyUserBeingFollowedRequest(ctx context.Context, event *model.UserFollowedEvent, req *model.NotifyUserBeingFollowedRequest) {
	req.FollowerID = event.FollowerID
	req.FollowingID = event.FollowingID
}
