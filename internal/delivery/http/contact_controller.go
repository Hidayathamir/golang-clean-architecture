package http

import (
	"errors"
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/delivery/http/response"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/pkg/httperror"
	"math"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ContactController struct {
	UseCase *usecase.ContactUseCase
	Log     *logrus.Logger
}

func NewContactController(useCase *usecase.ContactUseCase, log *logrus.Logger) *ContactController {
	return &ContactController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *ContactController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := new(model.CreateContactRequest)
	if err := ctx.BodyParser(req); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		err = errors.Join(httperror.BadRequest(), err)
		return response.Error(ctx, err)
	}
	req.UserId = auth.ID

	res, err := c.UseCase.Create(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

func (c *ContactController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := &model.SearchContactRequest{
		UserId: auth.ID,
		Name:   ctx.Query("name", ""),
		Email:  ctx.Query("email", ""),
		Phone:  ctx.Query("phone", ""),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	res, total, err := c.UseCase.Search(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("error searching contact")
		return response.Error(ctx, err)
	}

	paging := &response.PageMetadata{
		Page:      req.Page,
		Size:      req.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(req.Size))),
	}

	return response.DataPaging(ctx, http.StatusOK, res, paging)
}

func (c *ContactController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := &model.GetContactRequest{
		UserId: auth.ID,
		ID:     ctx.Params("contactId"),
	}

	res, err := c.UseCase.Get(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

func (c *ContactController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := new(model.UpdateContactRequest)
	if err := ctx.BodyParser(req); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		err = errors.Join(httperror.BadRequest(), err)
		return response.Error(ctx, err)
	}

	req.UserId = auth.ID
	req.ID = ctx.Params("contactId")

	res, err := c.UseCase.Update(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

func (c *ContactController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")

	req := &model.DeleteContactRequest{
		UserId: auth.ID,
		ID:     contactId,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), req); err != nil {
		c.Log.WithError(err).Error("error deleting contact")
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, true)
}
