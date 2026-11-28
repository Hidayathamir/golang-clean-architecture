package http

import (
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notifusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
)

type NotificationController struct {
	Cfg     *config.Config
	Usecase notifusecase.NotifUsecase
}

func NewNotificationController(cfg *config.Config, usecase notifusecase.NotifUsecase) *NotificationController {
	return &NotificationController{
		Cfg:     cfg,
		Usecase: usecase,
	}
}

// List godoc
//
//	@Summary		List notifications
//	@Description	List the current user's notifications, newest first
//	@Tags			notifications
//	@Param			limit	query		int		false	"page size"
//	@Param			cursor	query		string	false	"opaque cursor"
//	@Success		200		{object}	response.WebResponse[dto.NotificationResponseList]
//	@Router			/api/notifications [get]
func (c *NotificationController) List(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()
	req := dto.NotificationsRequest{}
	err := ctx.QueryParser(&req.Page)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "http.(*NotificationController).List")
	}
	res, err := c.Usecase.List(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*NotificationController).List")
	}
	return response.Data(ctx, http.StatusOK, res)
}

// MarkRead godoc
//
//	@Summary		Mark a notification read
//	@Tags			notifications
//	@Param			request	body		dto.MarkNotificationReadRequest	true	"Mark Read Request"
//	@Success		200		{object}	response.WebResponse[any]
//	@Router			/api/notifications/_read [post]
func (c *NotificationController) MarkRead(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()
	req := dto.MarkNotificationReadRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "http.(*NotificationController).MarkRead")
	}
	err = c.Usecase.MarkRead(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*NotificationController).MarkRead")
	}
	return response.Data(ctx, http.StatusOK, map[string]any{"id": req.NotificationID})
}

// MarkAllRead godoc
//
//	@Summary		Mark all notifications read
//	@Tags			notifications
//	@Success		200	{object}	response.WebResponse[any]
//	@Router			/api/notifications/_read_all [post]
func (c *NotificationController) MarkAllRead(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()
	err := c.Usecase.MarkAllRead(ctx.UserContext())
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*NotificationController).MarkAllRead")
	}
	return response.Data(ctx, http.StatusOK, map[string]any{"ok": true})
}

// UnreadCount godoc
//
//	@Summary		Unread notification count
//	@Tags			notifications
//	@Success		200	{object}	response.WebResponse[dto.UnreadCountResponse]
//	@Router			/api/notifications/_unread_count [get]
func (c *NotificationController) UnreadCount(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()
	res, err := c.Usecase.CountUnread(ctx.UserContext())
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*NotificationController).UnreadCount")
	}
	return response.Data(ctx, http.StatusOK, res)
}
