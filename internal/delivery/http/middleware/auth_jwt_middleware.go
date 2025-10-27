package middleware

import (
	"strings"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	user2usecase "github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user2"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/gofiber/fiber/v2"
)

func NewUser2Auth(user2Usecase user2usecase.User2Usecase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headerAuth := ctx.Get("Authorization")
		if headerAuth == "" {
			err := errkit.Unauthorized(fiber.ErrUnauthorized)
			return errkit.AddFuncName("middleware.NewUser2Auth", err)
		}

		token := strings.TrimSpace(headerAuth)
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = strings.TrimSpace(token[7:])
		}

		req := &model.VerifyUser2TokenRequest{Token: token}
		auth, err := user2Usecase.VerifyToken(ctx.UserContext(), req)
		if err != nil {
			err = errkit.Unauthorized(err)
			return errkit.AddFuncName("middleware.NewUser2Auth", err)
		}

		ctx.Locals("user2_auth", auth)
		return ctx.Next()
	}
}

func GetUser2(ctx *fiber.Ctx) *model.User2Auth {
	value := ctx.Locals("user2_auth")
	if value == nil {
		return nil
	}
	return value.(*model.User2Auth)
}
