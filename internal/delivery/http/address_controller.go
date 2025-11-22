package http

import (
	"net/http"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type AddressController struct {
	Config  *viper.Viper
	Usecase address.AddressUsecase
}

func NewAddressController(cfg *viper.Viper, useCase address.AddressUsecase) *AddressController {
	return &AddressController{
		Config:  cfg,
		Usecase: useCase,
	}
}

// Create godoc
//
//	@Summary		Create address
//	@Description	Create a new address for a contact
//	@Tags			addresses
//	@Param			contactId	path	string						true	"Contact ID"
//	@Param			request		body	model.CreateAddressRequest	true	"Create Address Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.AddressResponse]
//	@Router			/api/contacts/{contactId}/addresses [post]
func (c *AddressController) Create(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := new(model.CreateAddressRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req.UserID = auth.ID

	contactID, err := strconv.ParseInt(ctx.Params("contactId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}
	req.ContactID = contactID

	res, err := c.Usecase.Create(ctx.UserContext(), req)
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
//	@Param			contactId	path	string	true	"Contact ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.AddressResponseList]
//	@Router			/api/contacts/{contactId}/addresses [get]
func (c *AddressController) List(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	contactID, err := strconv.ParseInt(ctx.Params("contactId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req := &model.ListAddressRequest{
		UserID:    auth.ID,
		ContactID: contactID,
	}

	res, err := c.Usecase.List(ctx.UserContext(), req)
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
//	@Param			contactId	path	string	true	"Contact ID"
//	@Param			addressId	path	string	true	"Address ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.AddressResponse]
//	@Router			/api/contacts/{contactId}/addresses/{addressId} [get]
func (c *AddressController) Get(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	contactID, err := strconv.ParseInt(ctx.Params("contactId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	addressID, err := strconv.ParseInt(ctx.Params("addressId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req := &model.GetAddressRequest{
		UserID:    auth.ID,
		ContactID: contactID,
		ID:        addressID,
	}

	res, err := c.Usecase.Get(ctx.UserContext(), req)
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
//	@Param			contactId	path	string						true	"Contact ID"
//	@Param			addressId	path	string						true	"Address ID"
//	@Param			request		body	model.UpdateAddressRequest	true	"Update Address Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.AddressResponse]
//	@Router			/api/contacts/{contactId}/addresses/{addressId} [put]
func (c *AddressController) Update(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := new(model.UpdateAddressRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req.UserID = auth.ID

	contactID, err := strconv.ParseInt(ctx.Params("contactId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}
	req.ContactID = contactID

	addressID, err := strconv.ParseInt(ctx.Params("addressId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}
	req.ID = addressID

	res, err := c.Usecase.Update(ctx.UserContext(), req)
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
//	@Param			contactId	path	string	true	"Contact ID"
//	@Param			addressId	path	string	true	"Address ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[bool]
//	@Router			/api/contacts/{contactId}/addresses/{addressId} [delete]
func (c *AddressController) Delete(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	contactID, err := strconv.ParseInt(ctx.Params("contactId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	addressID, err := strconv.ParseInt(ctx.Params("addressId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req := &model.DeleteAddressRequest{
		UserID:    auth.ID,
		ContactID: contactID,
		ID:        addressID,
	}

	if err := c.Usecase.Delete(ctx.UserContext(), req); err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, true)
}
