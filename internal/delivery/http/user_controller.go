package http

import (
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/delivery/http/response"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/pkg/errkit"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UseCase usecase.UserUseCase
	Log     *logrus.Logger
}

func NewUserController(useCase usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		UseCase: useCase,
		Log:     logger,
	}
}

// Register godoc
//
//	@Summary		Register user
//	@Description	Register a new user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.RegisterUserRequest	true	"Register User Request"
//	@Success		200		{object}	response.WebResponse[model.UserResponse]
//	@Router			/api/users [post]
func (c *UserController) Register(ctx *fiber.Ctx) error {
	req := new(model.RegisterUserRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		err = errkit.BadRequest(err)
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	res, err := c.UseCase.Create(ctx.UserContext(), req)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate a user and return access token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.LoginUserRequest	true	"Login User Request"
//	@Success		200		{object}	response.WebResponse[model.UserResponse]
//	@Router			/api/users/_login [post]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	req := new(model.LoginUserRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		err = errkit.BadRequest(err)
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	res, err := c.UseCase.Login(ctx.UserContext(), req)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Current godoc
//
//	@Summary		Get current user
//	@Description	Get profile of the currently authenticated user
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.WebResponse[model.UserResponse]
//	@Router			/api/users/_current [get]
func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := &model.GetUserRequest{
		ID: auth.ID,
	}

	res, err := c.UseCase.Current(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current user")
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Logout godoc
//
//	@Summary		Logout user
//	@Description	Logout the current authenticated user
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.WebResponse[bool]
//	@Router			/api/users [delete]
func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := &model.LogoutUserRequest{
		ID: auth.ID,
	}

	res, err := c.UseCase.Logout(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to logout user")
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Update godoc
//
//	@Summary		Update user
//	@Description	Update profile of the current authenticated user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		model.UpdateUserRequest	true	"Update User Request"
//	@Success		200		{object}	response.WebResponse[model.UserResponse]
//	@Router			/api/users/_current [patch]
func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(req); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		err = errkit.BadRequest(err)
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	req.ID = auth.ID
	res, err := c.UseCase.Update(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to update user")
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}
