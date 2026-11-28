package http

import (
	"net/http"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/chatusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
)

type ChatController struct {
	Cfg     *config.Config
	Usecase chatusecase.ChatUsecase
}

func NewChatController(cfg *config.Config, usecase chatusecase.ChatUsecase) *ChatController {
	return &ChatController{
		Cfg:     cfg,
		Usecase: usecase,
	}
}

// Send godoc
//
//	@Summary	Send a chat
//	@Tags		chats
//	@Param		request	body		dto.SendChatRequest	true	"Send Chat Request"
//	@Success	200		{object}	response.WebResponse[dto.ChatResponse]
//	@Router		/api/chats [post]
func (c *ChatController) Send(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()
	req := dto.SendChatRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "http.(*ChatController).Send")
	}
	res, err := c.Usecase.Send(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ChatController).Send")
	}
	return response.Data(ctx, http.StatusOK, res)
}

// ListConversations godoc
//
//	@Summary	List chat conversations
//	@Tags		chats
//	@Success	200	{object}	response.WebResponse[dto.ConversationsResponse]
//	@Router		/api/chats/conversations [get]
func (c *ChatController) ListConversations(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()
	res, err := c.Usecase.ListConversations(ctx.UserContext())
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ChatController).ListConversations")
	}
	return response.Data(ctx, http.StatusOK, res)
}

// GetThread godoc
//
//	@Summary	Get a chat thread with a user
//	@Tags		chats
//	@Param		userId	path		int		true	"peer user id"
//	@Param		limit	query		int		false	"page size"
//	@Param		cursor	query		string	false	"opaque cursor"
//	@Success	200		{object}	response.WebResponse[dto.ThreadResponse]
//	@Router		/api/chats/{userId} [get]
func (c *ChatController) GetThread(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()
	req := dto.GetThreadRequest{}
	peerID, err := strconv.ParseInt(ctx.Params("userId"), 10, 64)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "http.(*ChatController).GetThread")
	}
	req.PeerID = peerID
	err = ctx.QueryParser(&req.Page)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "http.(*ChatController).GetThread")
	}
	res, err := c.Usecase.GetThread(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ChatController).GetThread")
	}
	return response.Data(ctx, http.StatusOK, res)
}

// MarkRead godoc
//
//	@Summary	Mark a chat thread read
//	@Tags		chats
//	@Param		request	body		dto.MarkChatReadRequest	true	"Mark Chat Read Request"
//	@Success	200		{object}	response.WebResponse[any]
//	@Router		/api/chats/_read [post]
func (c *ChatController) MarkRead(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()
	req := dto.MarkChatReadRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "http.(*ChatController).MarkRead")
	}
	err = c.Usecase.MarkRead(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ChatController).MarkRead")
	}
	return response.Data(ctx, http.StatusOK, map[string]any{"ok": true})
}

// UnreadCount godoc
//
//	@Summary	Unread chat count
//	@Tags		chats
//	@Success	200	{object}	response.WebResponse[dto.UnreadChatCountResponse]
//	@Router		/api/chats/_unread_count [get]
func (c *ChatController) UnreadCount(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()
	res, err := c.Usecase.CountUnread(ctx.UserContext())
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ChatController).UnreadCount")
	}
	return response.Data(ctx, http.StatusOK, res)
}
