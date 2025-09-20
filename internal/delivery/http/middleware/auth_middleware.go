package middleware

import (
	"golang-clean-architecture/internal/delivery/http/response"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/pkg/errkit"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewAuth(log *logrus.Logger, userUserCase usecase.UserUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		log.Debugf("Authorization : %s", request.Token)

		auth, err := userUserCase.Verify(ctx.UserContext(), request)
		if err != nil {
			log.Warnf("Failed find user by token : %+v", err)
			err = errkit.Unauthorized(err)
			err = errkit.AddFuncName(err)
			return response.Error(ctx, err)
		}

		log.Debugf("User : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
