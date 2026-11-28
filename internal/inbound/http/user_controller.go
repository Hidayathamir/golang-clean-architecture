package http

import (
	"net/http"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Cfg     *config.Config
	Usecase userusecase.UserUsecase
}

func NewUserController(cfg *config.Config, usecase userusecase.UserUsecase) *UserController {
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
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Register")
	}

	res, err := c.Usecase.Create(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Register")
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
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Login")
	}

	res, err := c.Usecase.Login(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Login")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// VerifyEmail godoc
//
//	@Summary		Verify email
//	@Description	Confirm a user's email address using the token from the verification email
//	@Tags			users
//	@Param			token	query		string	true	"Verification token"
//	@Success		200		{object}	response.WebResponse[string]
//	@Router			/api/users/_verify_email [get]
func (c *UserController) VerifyEmail(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := dto.VerifyEmailRequest{
		Token: ctx.Query("token"),
	}

	err := c.Usecase.VerifyEmail(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).VerifyEmail")
	}

	return response.Data(ctx, http.StatusOK, "ok")
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
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Current")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Get a user's public profile with follower/following counts
//	@Tags			users
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.ProfileResponse]
//	@Router			/api/users/{id} [get]
func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	userID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).GetProfile")
	}

	req := dto.GetProfileRequest{
		ID: userID,
	}

	res, err := c.Usecase.GetProfile(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).GetProfile")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// GetFollowers godoc
//
//	@Summary		Get user followers
//	@Description	Get the list of users following the given user (cursor paginated)
//	@Tags			users
//	@Produce		json
//	@Param			id		path	int		true	"User ID"
//	@Param			limit	query	int		false	"Page size (default 20, max 100)"
//	@Param			cursor	query	string	false	"Opaque cursor from the previous page; empty for the first page"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.FollowersResponse]
//	@Router			/api/users/{id}/followers [get]
func (c *UserController) GetFollowers(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	userID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).GetFollowers")
	}

	req := dto.FollowersRequest{
		UserID: userID,
	}
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			err = errkit.SetCode(err, http.StatusBadRequest)
			logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
			return errkit.AddFuncName(err, "http.(*UserController).GetFollowers")
		}
		req.Page.Limit = limit
	}
	req.Page.Cursor = ctx.Query("cursor")

	res, err := c.Usecase.GetFollowers(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).GetFollowers")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// GetFollowings godoc
//
//	@Summary		Get user followings
//	@Description	Get the list of users the given user follows (cursor paginated)
//	@Tags			users
//	@Produce		json
//	@Param			id		path	int		true	"User ID"
//	@Param			limit	query	int		false	"Page size (default 20, max 100)"
//	@Param			cursor	query	string	false	"Opaque cursor from the previous page; empty for the first page"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.FollowingsResponse]
//	@Router			/api/users/{id}/following [get]
func (c *UserController) GetFollowings(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	userID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).GetFollowings")
	}

	req := dto.FollowingsRequest{
		UserID: userID,
	}
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			err = errkit.SetCode(err, http.StatusBadRequest)
			logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
			return errkit.AddFuncName(err, "http.(*UserController).GetFollowings")
		}
		req.Page.Limit = limit
	}
	req.Page.Cursor = ctx.Query("cursor")

	res, err := c.Usecase.GetFollowings(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).GetFollowings")
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Unfollow godoc
//
//	@Summary		Unfollow user
//	@Description	Unfollow a user the current user follows
//	@Tags			users
//	@Produce		json
//	@Param			id	path	int	true	"User ID to unfollow"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[string]
//	@Router			/api/users/{id}/follow [delete]
func (c *UserController) Unfollow(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	followingID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Unfollow")
	}

	req := dto.UnfollowUserRequest{
		FollowingID: followingID,
	}

	err = c.Usecase.Unfollow(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Unfollow")
	}

	return response.Data(ctx, http.StatusOK, "ok")
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
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Update")
	}

	req.ID = userAuth.ID
	res, err := c.Usecase.Update(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Update")
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
		err = errkit.SetCode(err, http.StatusBadRequest)
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Follow")
	}

	err = c.Usecase.Follow(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*UserController).Follow")
	}

	return response.Data(ctx, http.StatusOK, "ok")
}
