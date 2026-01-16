package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
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
	if user == nil || res == nil {
		return
	}

	res.ID = user.ID
	res.Username = user.Username
	res.Name = user.Name
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
}

func EntityUserToModelUserFollowedEvent(user *entity.User, event *model.UserFollowedEvent) {
	if user == nil || event == nil {
		return
	}

	event.ID = user.ID
	event.Username = user.Username
	event.Name = user.Name
	event.CreatedAt = user.CreatedAt
	event.UpdatedAt = user.UpdatedAt
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
