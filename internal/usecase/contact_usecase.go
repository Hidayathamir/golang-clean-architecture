package usecase

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/pkg/errkit"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ContactUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	ContactRepository *repository.ContactRepository
	ContactProducer   *messaging.ContactProducer
}

func NewContactUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	contactRepository *repository.ContactRepository, contactProducer *messaging.ContactProducer) *ContactUseCase {
	return &ContactUseCase{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		ContactRepository: contactRepository,
		ContactProducer:   contactProducer,
	}
}

func (c *ContactUseCase) Create(ctx context.Context, request *model.CreateContactRequest) (*model.ContactResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, errkit.BadRequest(err)
	}

	contact := &entity.Contact{
		ID:        uuid.New().String(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Phone:     request.Phone,
		UserId:    request.UserId,
	}

	if err := c.ContactRepository.Create(tx, contact); err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return nil, errkit.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating contact")
		return nil, errkit.InternalServerError(err)
	}

	event := converter.ContactToEvent(contact)
	if err := c.ContactProducer.Send(event); err != nil {
		c.Log.WithError(err).Error("error publishing contact created event")
		return nil, errkit.InternalServerError(err)
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Update(ctx context.Context, request *model.UpdateContactRequest) (*model.ContactResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ID, request.UserId); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, errkit.NotFound(err)
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, errkit.BadRequest(err)
	}

	contact.FirstName = request.FirstName
	contact.LastName = request.LastName
	contact.Email = request.Email
	contact.Phone = request.Phone

	if err := c.ContactRepository.Update(tx, contact); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, errkit.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, errkit.InternalServerError(err)
	}

	event := converter.ContactToEvent(contact)
	if err := c.ContactProducer.Send(event); err != nil {
		c.Log.WithError(err).Error("error publishing contact updated event")
		return nil, errkit.InternalServerError(err)
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Get(ctx context.Context, request *model.GetContactRequest) (*model.ContactResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, errkit.BadRequest(err)
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ID, request.UserId); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, errkit.NotFound(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, errkit.InternalServerError(err)
	}

	return converter.ContactToResponse(contact), nil
}

func (c *ContactUseCase) Delete(ctx context.Context, request *model.DeleteContactRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return errkit.BadRequest(err)
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ID, request.UserId); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return errkit.NotFound(err)
	}

	if err := c.ContactRepository.Delete(tx, contact); err != nil {
		c.Log.WithError(err).Error("error deleting contact")
		return errkit.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error deleting contact")
		return errkit.InternalServerError(err)
	}

	return nil
}

func (c *ContactUseCase) Search(ctx context.Context, request *model.SearchContactRequest) ([]model.ContactResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, 0, errkit.BadRequest(err)
	}

	contacts, total, err := c.ContactRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting contacts")
		return nil, 0, errkit.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error getting contacts")
		return nil, 0, errkit.InternalServerError(err)
	}

	responses := make([]model.ContactResponse, len(contacts))
	for i, contact := range contacts {
		responses[i] = *converter.ContactToResponse(&contact)
	}

	return responses, total, nil
}
