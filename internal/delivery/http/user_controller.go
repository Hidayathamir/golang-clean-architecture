package http

import (
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Cfg     *config.Config
	Usecase user.UserUsecase
}

func NewUserController(cfg *config.Config, usecase user.UserUsecase) *UserController {
	return &UserController{
		Cfg:     cfg,
		Usecase: usecase,
	}
}

// Register godoc
//
//	@Summary		Register user
//	@Description	Register a new user account
//	@Tags			users
//	@Param			request	body		dto.RegisterUserRequest	true	"Register User Request"
//	@Success		200		{object}	response.WebResponse[dto.UserResponse]
//	@Router			/api/users [post]
func (c *UserController) Register(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := dto.RegisterUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	res, err := c.Usecase.Create(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate a user and return access token
//	@Tags			users
//	@Param			request	body		dto.LoginUserRequest	true	"Login User Request"
//	@Success		200		{object}	response.WebResponse[dto.UserLoginResponse]
//	@Router			/api/users/_login [post]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := dto.LoginUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	res, err := c.Usecase.Login(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Current godoc
//
//	@Summary		Get current user
//	@Description	Get profile of the currently authenticated user
//	@Tags			users
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.UserResponse]
//	@Router			/api/users/_current [get]
func (c *UserController) Current(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	userAuth := ctxuserauth.Get(ctx.UserContext())

	req := dto.GetUserRequest{
		ID: userAuth.ID,
	}

	res, err := c.Usecase.Current(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Update godoc
//
//	@Summary		Update user
//	@Description	Update profile of the current authenticated user
//	@Tags			users
//	@Security		SimpleApiKeyAuth
//	@Param			request	body		dto.UpdateUserRequest	true	"Update User Request"
//	@Success		200		{object}	response.WebResponse[dto.UserResponse]
//	@Router			/api/users/_current [patch]
func (c *UserController) Update(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	userAuth := ctxuserauth.Get(ctx.UserContext())

	req := dto.UpdateUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req.ID = userAuth.ID
	res, err := c.Usecase.Update(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Follow godoc
//
//	@Summary		Follow user
//	@Description	Follow a user
//	@Tags			users
//	@Security		SimpleApiKeyAuth
//	@Param			request	body		dto.FollowUserRequest	true	"Follow User Request"
//	@Success		200		{object}	response.WebResponse[string]
//	@Router			/api/users/_follow [post]
func (c *UserController) Follow(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := dto.FollowUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	err = c.Usecase.Follow(ctx.UserContext(), req)
	if err != nil {
		x.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, "ok")
}
