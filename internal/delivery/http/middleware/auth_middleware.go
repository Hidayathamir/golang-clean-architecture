package middleware

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUserCase user.UserUsecase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}

		auth, err := userUserCase.Verify(ctx.UserContext(), request)
		if err != nil {
			err = errkit.Unauthorized(err)
			return errkit.AddFuncName(err)
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
