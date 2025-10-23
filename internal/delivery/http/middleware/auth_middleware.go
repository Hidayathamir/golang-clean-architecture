package middleware

import (
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUserCase user.UserUsecase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headerAuth := ctx.Get("Authorization")
		if headerAuth == "" {
			err := fmt.Errorf("header auth not found")
			err = errkit.Unauthorized(err)
			return errkit.AddFuncName("middleware.NewAuth", err)
		}

		req := &model.VerifyUserRequest{Token: headerAuth}
		auth, err := userUserCase.Verify(ctx.UserContext(), req)
		if err != nil {
			err = errkit.Unauthorized(err)
			return errkit.AddFuncName("middleware.NewAuth", err)
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
