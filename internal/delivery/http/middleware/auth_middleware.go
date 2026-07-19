package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUserCase user.UserUsecase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headerAuth := strings.TrimSpace(ctx.Get("Authorization"))
		if headerAuth == "" {
			err := fmt.Errorf("header auth not found")
			err = errkit.SetCode(err, http.StatusUnauthorized)
			return errkit.AddFuncName(err)
		}

		var token string
		parts := strings.Fields(headerAuth)
		switch {
		case len(parts) == 1:
			token = parts[0]
		case len(parts) == 2 && strings.EqualFold(parts[0], "Bearer"):
			token = parts[1]
		default:
			err := fmt.Errorf("authorization header format invalid")
			err = errkit.SetCode(err, http.StatusUnauthorized)
			return errkit.AddFuncName(err)
		}

		req := dto.VerifyUserRequest{Token: token}
		userAuth, err := userUserCase.Verify(ctx.UserContext(), req)
		if err != nil {
			err = errkit.SetCode(err, http.StatusUnauthorized)
			return errkit.AddFuncName(err)
		}

		ctx.SetUserContext(ctxuserauth.Set(ctx.UserContext(), &userAuth))

		return ctx.Next()
	}
}
