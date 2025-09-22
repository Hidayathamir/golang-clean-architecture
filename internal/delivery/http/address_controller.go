package http

import (
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/delivery/http/response"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase/address"
	"golang-clean-architecture/pkg/errkit"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AddressController struct {
	UseCase address.AddressUseCase
	Log     *logrus.Logger
}

func NewAddressController(useCase address.AddressUseCase, log *logrus.Logger) *AddressController {
	return &AddressController{
		UseCase: useCase,
		Log:     log,
	}
}

// Create godoc
//
//	@Summary		Create address
//	@Description	Create a new address for a contact
//	@Tags			addresses
//	@Accept			json
//	@Produce		json
//	@Param			contactId	path		string						true	"Contact ID"
//	@Param			request		body		model.CreateAddressRequest	true	"Create Address Request"
//	@Success		200			{object}	response.WebResponse[model.AddressResponse]
//	@Router			/api/contacts/{contactId}/addresses [post]
func (c *AddressController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := new(model.CreateAddressRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req.UserId = auth.ID
	req.ContactId = ctx.Params("contactId")

	res, err := c.UseCase.Create(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// List godoc
//
//	@Summary		List addresses
//	@Description	Get all addresses for a contact
//	@Tags			addresses
//	@Produce		json
//	@Param			contactId	path		string	true	"Contact ID"
//	@Success		200			{object}	response.WebResponse[[]model.AddressResponse]
//	@Router			/api/contacts/{contactId}/addresses [get]
func (c *AddressController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := ctx.Params("contactId")

	req := &model.ListAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
	}

	res, err := c.UseCase.List(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Get godoc
//
//	@Summary		Get address
//	@Description	Get a specific address by ID
//	@Tags			addresses
//	@Produce		json
//	@Param			contactId	path		string	true	"Contact ID"
//	@Param			addressId	path		string	true	"Address ID"
//	@Success		200			{object}	response.WebResponse[model.AddressResponse]
//	@Router			/api/contacts/{contactId}/addresses/{addressId} [get]
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
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Update godoc
//
//	@Summary		Update address
//	@Description	Update an existing address by ID
//	@Tags			addresses
//	@Accept			json
//	@Produce		json
//	@Param			contactId	path		string						true	"Contact ID"
//	@Param			addressId	path		string						true	"Address ID"
//	@Param			request		body		model.UpdateAddressRequest	true	"Update Address Request"
//	@Success		200			{object}	response.WebResponse[model.AddressResponse]
//	@Router			/api/contacts/{contactId}/addresses/{addressId} [put]
func (c *AddressController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	req := new(model.UpdateAddressRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req.UserId = auth.ID
	req.ContactId = ctx.Params("contactId")
	req.ID = ctx.Params("addressId")

	res, err := c.UseCase.Update(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Delete godoc
//
//	@Summary		Delete address
//	@Description	Delete an address by ID
//	@Tags			addresses
//	@Produce		json
//	@Param			contactId	path		string	true	"Contact ID"
//	@Param			addressId	path		string	true	"Address ID"
//	@Success		200			{object}	response.WebResponse[bool]
//	@Router			/api/contacts/{contactId}/addresses/{addressId} [delete]
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
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, true)
}
