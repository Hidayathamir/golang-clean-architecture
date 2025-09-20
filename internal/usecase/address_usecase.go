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

type AddressUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	AddressRepository *repository.AddressRepository
	ContactRepository *repository.ContactRepository
	AddressProducer   *messaging.AddressProducer
}

func NewAddressUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	contactRepository *repository.ContactRepository, addressRepository *repository.AddressRepository,
	addressProducer *messaging.AddressProducer) *AddressUseCase {
	return &AddressUseCase{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		ContactRepository: contactRepository,
		AddressRepository: addressRepository,
		AddressProducer:   addressProducer,
	}
}

func (c *AddressUseCase) Create(ctx context.Context, request *model.CreateAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, errkit.BadRequest(err)
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, errkit.NotFound(err)
	}

	address := &entity.Address{
		ID:         uuid.NewString(),
		ContactId:  contact.ID,
		Street:     request.Street,
		City:       request.City,
		Province:   request.Province,
		PostalCode: request.PostalCode,
		Country:    request.Country,
	}

	if err := c.AddressRepository.Create(tx, address); err != nil {
		c.Log.WithError(err).Error("failed to create address")
		return nil, errkit.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, errkit.InternalServerError(err)
	}

	if c.AddressProducer != nil {
		event := converter.AddressToEvent(address)
		if err := c.AddressProducer.Send(event); err != nil {
			c.Log.WithError(err).Error("failed to publish address created event")
			return nil, errkit.InternalServerError(err)
		}
		c.Log.Info("Published address created event")
	} else {
		c.Log.Info("Kafka producer is disabled, skipping address created event")
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Update(ctx context.Context, request *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, errkit.BadRequest(err)
	}

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, errkit.NotFound(err)
	}

	address := new(entity.Address)
	if err := c.AddressRepository.FindByIdAndContactId(tx, address, request.ID, contact.ID); err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return nil, errkit.NotFound(err)
	}

	address.Street = request.Street
	address.City = request.City
	address.Province = request.Province
	address.PostalCode = request.PostalCode
	address.Country = request.Country

	if err := c.AddressRepository.Update(tx, address); err != nil {
		c.Log.WithError(err).Error("failed to update address")
		return nil, errkit.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, errkit.InternalServerError(err)
	}

	event := converter.AddressToEvent(address)
	if err := c.AddressProducer.Send(event); err != nil {
		c.Log.WithError(err).Error("failed to publish address updated event")
		return nil, errkit.InternalServerError(err)
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Get(ctx context.Context, request *model.GetAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, errkit.NotFound(err)
	}

	address := new(entity.Address)
	if err := c.AddressRepository.FindByIdAndContactId(tx, address, request.ID, request.ContactId); err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return nil, errkit.NotFound(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, errkit.InternalServerError(err)
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) Delete(ctx context.Context, request *model.DeleteAddressRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return errkit.NotFound(err)
	}

	address := new(entity.Address)
	if err := c.AddressRepository.FindByIdAndContactId(tx, address, request.ID, request.ContactId); err != nil {
		c.Log.WithError(err).Error("failed to find address")
		return errkit.NotFound(err)
	}

	if err := c.AddressRepository.Delete(tx, address); err != nil {
		c.Log.WithError(err).Error("failed to delete address")
		return errkit.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return errkit.InternalServerError(err)
	}

	return nil
}

func (c *AddressUseCase) List(ctx context.Context, request *model.ListAddressRequest) ([]model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := c.ContactRepository.FindByIdAndUserId(tx, contact, request.ContactId, request.UserId); err != nil {
		c.Log.WithError(err).Error("failed to find contact")
		return nil, errkit.NotFound(err)
	}

	addresses, err := c.AddressRepository.FindAllByContactId(tx, contact.ID)
	if err != nil {
		c.Log.WithError(err).Error("failed to find addresses")
		return nil, errkit.InternalServerError(err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, errkit.InternalServerError(err)
	}

	responses := make([]model.AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = *converter.AddressToResponse(&address)
	}

	return responses, nil
}
