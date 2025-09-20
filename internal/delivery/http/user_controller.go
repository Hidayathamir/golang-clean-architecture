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
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

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
