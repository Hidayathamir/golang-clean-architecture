package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
)

func ModelRegisterUser2RequestToEntityUser2(req *model.RegisterUser2Request, user *entity.User2, hashedPassword string, id string) {
	if req == nil || user == nil {
		return
	}

	user.ID = id
	user.Email = req.Email
	user.DisplayName = req.DisplayName
	user.Password = hashedPassword
}

func EntityUser2ToModelUser2Response(user *entity.User2, res *model.User2Response) {
	if user == nil || res == nil {
		return
	}

	res.ID = user.ID
	res.Email = user.Email
	res.DisplayName = user.DisplayName
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
}
