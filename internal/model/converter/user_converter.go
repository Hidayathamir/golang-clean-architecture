package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
)

func ModelRegisterUserRequestToEntityUser(req *model.RegisterUserRequest, user *entity.User, password string) {
	if req == nil || user == nil {
		return
	}

	user.ID = req.ID
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

func UserToResponse(user *entity.User, res *model.UserResponse) {
	if user == nil || res == nil {
		return
	}

	res.ID = user.ID
	res.Name = user.Name
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
}

func UserToTokenResponse(user *entity.User, res *model.UserResponse) {
	if user == nil || res == nil {
		return
	}

	res.Token = user.Token
}

func UserToEvent(user *entity.User, event *model.UserEvent) {
	if user == nil || event == nil {
		return
	}

	event.ID = user.ID
	event.Name = user.Name
	event.CreatedAt = user.CreatedAt
	event.UpdatedAt = user.UpdatedAt
}
