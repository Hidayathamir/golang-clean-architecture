package http

import (
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UserController struct {
	Config  *viper.Viper
	Log     *logrus.Logger
	Usecase user.UserUsecase
}

func NewUserController(cfg *viper.Viper, log *logrus.Logger, useCase user.UserUsecase) *UserController {
	return &UserController{
		Config:  cfg,
		Log:     log,
		Usecase: useCase,
	}
}

// Register godoc
//
//	@Summary		Register user
//	@Description	Register a new user account
//	@Tags			users
//	@Param			request	body		model.RegisterUserRequest	true	"Register User Request"
//	@Success		200		{object}	response.WebResponse[model.UserResponse]
//	@Router			/api/users [post]
func (c *UserController) Register(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := new(model.RegisterUserRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName("http.(*UserController).Register", err)
	}

	res, err := c.Usecase.Create(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*UserController).Register", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate a user and return access token
//	@Tags			users
//	@Param			request	body		model.LoginUserRequest	true	"Login User Request"
//	@Success		200		{object}	response.WebResponse[model.UserResponse]
//	@Router			/api/users/_login [post]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := new(model.LoginUserRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName("http.(*UserController).Login", err)
	}

	res, err := c.Usecase.Login(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*UserController).Login", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Current godoc
//
//	@Summary		Get current user
//	@Description	Get profile of the currently authenticated user
//	@Tags			users
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.UserResponse]
//	@Router			/api/users/_current [get]
func (c *UserController) Current(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := &model.GetUserRequest{
		ID: auth.ID,
	}

	res, err := c.Usecase.Current(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*UserController).Current", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Logout godoc
//
//	@Summary		Logout user
//	@Description	Logout the current authenticated user
//	@Tags			users
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[bool]
//	@Router			/api/users [delete]
func (c *UserController) Logout(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := &model.LogoutUserRequest{
		ID: auth.ID,
	}

	res, err := c.Usecase.Logout(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*UserController).Logout", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Update godoc
//
//	@Summary		Update user
//	@Description	Update profile of the current authenticated user
//	@Tags			users
//	@Security		SimpleApiKeyAuth
//	@Param			request	body		model.UpdateUserRequest	true	"Update User Request"
//	@Success		200		{object}	response.WebResponse[model.UserResponse]
//	@Router			/api/users/_current [patch]
func (c *UserController) Update(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName("http.(*UserController).Update", err)
	}

	req.ID = auth.ID
	res, err := c.Usecase.Update(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName("http.(*UserController).Update", err)
	}

	return response.Data(ctx, http.StatusOK, res)
}
