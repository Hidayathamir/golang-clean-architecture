package http

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type ContactController struct {
	Config  *viper.Viper
	Usecase contact.ContactUsecase
}

func NewContactController(cfg *viper.Viper, useCase contact.ContactUsecase) *ContactController {
	return &ContactController{
		Config:  cfg,
		Usecase: useCase,
	}
}

// Create godoc
//
//	@Summary		Create contact
//	@Description	Create a new contact
//	@Tags			contacts
//	@Param			request	body	model.CreateContactRequest	true	"Create Contact Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.ContactResponse]
//	@Router			/api/contacts [post]
func (c *ContactController) Create(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := new(model.CreateContactRequest)
	if err := ctx.BodyParser(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}
	req.UserID = auth.ID

	res, err := c.Usecase.Create(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// List godoc
//
//	@Summary		List contacts
//	@Description	Search and list contacts with filters and pagination
//	@Tags			contacts
//	@Param			name	query	string	false	"Filter by name"
//	@Param			email	query	string	false	"Filter by email"
//	@Param			phone	query	string	false	"Filter by phone"
//	@Param			page	query	int		false	"Page number"	default(1)
//	@Param			size	query	int		false	"Page size"		default(10)
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.ContactResponseList]
//	@Router			/api/contacts [get]
func (c *ContactController) List(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := &model.SearchContactRequest{
		UserID: auth.ID,
		Name:   ctx.Query("name", ""),
		Email:  ctx.Query("email", ""),
		Phone:  ctx.Query("phone", ""),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	res, total, err := c.Usecase.Search(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	paging := &response.PageMetadata{
		Page:      req.Page,
		Size:      req.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(req.Size))),
	}

	return response.DataPaging(ctx, http.StatusOK, res, paging)
}

// Get godoc
//
//	@Summary		Get contact
//	@Description	Get a specific contact by ID
//	@Tags			contacts
//	@Param			contactId	path	string	true	"Contact ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.ContactResponse]
//	@Router			/api/contacts/{contactId} [get]
func (c *ContactController) Get(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	contactID, err := strconv.ParseInt(ctx.Params("contactId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req := &model.GetContactRequest{
		UserID: auth.ID,
		ID:     contactID,
	}

	res, err := c.Usecase.Get(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Update godoc
//
//	@Summary		Update contact
//	@Description	Update an existing contact by ID
//	@Tags			contacts
//	@Param			contactId	path	string						true	"Contact ID"
//	@Param			request		body	model.UpdateContactRequest	true	"Update Contact Request"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[model.ContactResponse]
//	@Router			/api/contacts/{contactId} [put]
func (c *ContactController) Update(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	req := new(model.UpdateContactRequest)
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
	req.ID = contactID

	res, err := c.Usecase.Update(ctx.UserContext(), req)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, res)
}

// Delete godoc
//
//	@Summary		Delete contact
//	@Description	Delete a contact by ID
//	@Tags			contacts
//	@Param			contactId	path	string	true	"Contact ID"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[bool]
//	@Router			/api/contacts/{contactId} [delete]
func (c *ContactController) Delete(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	auth := middleware.GetUser(ctx)

	contactID, err := strconv.ParseInt(ctx.Params("contactId"), 10, 64)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	req := &model.DeleteContactRequest{
		UserID: auth.ID,
		ID:     contactID,
	}

	if err := c.Usecase.Delete(ctx.UserContext(), req); err != nil {
		return errkit.AddFuncName(err)
	}

	return response.Data(ctx, http.StatusOK, true)
}
