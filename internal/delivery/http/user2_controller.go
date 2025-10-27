package http

import (
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	user2usecase "github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user2"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type User2Controller struct {
	Config  *viper.Viper
	Log     *logrus.Logger
	Usecase user2usecase.User2Usecase
}

func NewUser2Controller(cfg *viper.Viper, log *logrus.Logger, usecase user2usecase.User2Usecase) *User2Controller {
	return &User2Controller{
		Config:  cfg,
		Log:     log,
		Usecase: usecase,
	}
}

// Register godoc
//
//	@Summary		Register user2 account
//	@Description	Create a new user2 account and return its profile
//	@Tags			user2
//	@Param			request	body		model.RegisterUser2Request	true	"Register User2 Request"
//	@Success		200		{object}	response.WebResponse[model.User2Response]
//	@Router			/api/user2 [post]
func (c *User2Controller) Register(ctx *fiber.Ctx) error {
	req := new(model.RegisterUser2Request)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName("http.(*User2Controller).Register", err)
	}

	res, err := c.Usecase.Register(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*User2Controller).Register", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Login godoc
//
//	@Summary		JWT login for user2
//	@Description	Authenticate a user2 account and return an access token
//	@Tags			user2
//	@Param			request	body		model.LoginUser2Request	true	"Login User2 Request"
//	@Success		200		{object}	response.WebResponse[model.User2TokenResponse]
//	@Router			/api/user2/_login [post]
func (c *User2Controller) Login(ctx *fiber.Ctx) error {
	req := new(model.LoginUser2Request)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName("http.(*User2Controller).Login", err)
	}

	res, err := c.Usecase.Login(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*User2Controller).Login", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Current godoc
//
//	@Summary		Get current user2 profile
//	@Description	Return the profile associated with the provided JWT access token
//	@Tags			user2
//	@Security		BearerAuth
//	@Success		200	{object}	response.WebResponse[model.User2Response]
//	@Router			/api/user2/_current [get]
func (c *User2Controller) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser2(ctx)
	if auth == nil {
		err := errkit.Unauthorized(fiber.ErrUnauthorized)
		return errkit.AddFuncName("http.(*User2Controller).Current", err)
	}

	res, err := c.Usecase.Profile(ctx.UserContext(), &model.GetUser2Request{ID: auth.UserID})
	if err != nil {
		return errkit.AddFuncName("http.(*User2Controller).Current", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}
