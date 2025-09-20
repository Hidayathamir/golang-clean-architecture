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

type AddressController struct {
	UseCase usecase.AddressUseCase
	Log     *logrus.Logger
}

func NewAddressController(useCase usecase.AddressUseCase, log *logrus.Logger) *AddressController {
	return &AddressController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *AddressController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := new(model.CreateAddressRequest)
	if err := ctx.BodyParser(req); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		err = errkit.BadRequest(err)
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	req.UserId = auth.ID
	req.ContactId = ctx.Params("contactId")

	res, err := c.UseCase.Create(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("failed to create address")
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

func (c *AddressController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")

	req := &model.ListAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
	}

	res, err := c.UseCase.List(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("failed to list addresses")
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

func (c *AddressController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")
	addressId := ctx.Params("addressId")

	req := &model.GetAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	res, err := c.UseCase.Get(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("failed to get address")
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

func (c *AddressController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := new(model.UpdateAddressRequest)
	if err := ctx.BodyParser(req); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		err = errkit.BadRequest(err)
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	req.UserId = auth.ID
	req.ContactId = ctx.Params("contactId")
	req.ID = ctx.Params("addressId")

	res, err := c.UseCase.Update(ctx.UserContext(), req)
	if err != nil {
		c.Log.WithError(err).Error("failed to update address")
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

func (c *AddressController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")
	addressId := ctx.Params("addressId")

	req := &model.DeleteAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
		ID:        addressId,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), req); err != nil {
		c.Log.WithError(err).Error("failed to delete address")
		err = errkit.AddFuncName(err)
		return response.Error(ctx, err)
	}

	return response.Data(ctx, http.StatusOK, true)
}
