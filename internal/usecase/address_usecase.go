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

type AddressUseCase interface {
	Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error)
	Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error)
	Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error)
	Delete(ctx context.Context, req *model.DeleteAddressRequest) error
	List(ctx context.Context, req *model.ListAddressRequest) ([]model.AddressResponse, error)
}

var _ AddressUseCase = &AddressUseCaseImpl{}

type AddressUseCaseImpl struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate

	// repository
	AddressRepository repository.AddressRepository
	ContactRepository repository.ContactRepository

	// producer
	AddressProducer messaging.AddressProducer
}

func NewAddressUseCase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,

	// repository
	contactRepository repository.ContactRepository,
	addressRepository repository.AddressRepository,

	// producer
	addressProducer messaging.AddressProducer,
) *AddressUseCaseImpl {
	return &AddressUseCaseImpl{
		DB:       db,
		Log:      logger,
		Validate: validate,

		// repository
		ContactRepository: contactRepository,
		AddressRepository: addressRepository,

		// producer
		AddressProducer: addressProducer,
	}
}

func (u *AddressUseCaseImpl) Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		u.Log.WithError(err).Error("failed to validate request body")
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, req.ContactId, req.UserId); err != nil {
		u.Log.WithError(err).Error("failed to find contact")
		return nil, errkit.AddFuncName(err)
	}

	address := &entity.Address{
		ID:         uuid.NewString(),
		ContactId:  contact.ID,
		Street:     req.Street,
		City:       req.City,
		Province:   req.Province,
		PostalCode: req.PostalCode,
		Country:    req.Country,
	}

	if err := u.AddressRepository.Create(tx, address); err != nil {
		u.Log.WithError(err).Error("failed to create address")
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("failed to commit transaction")
		return nil, errkit.AddFuncName(err)
	}

	if u.AddressProducer != nil {
		event := converter.AddressToEvent(address)
		if err := u.AddressProducer.Send(event); err != nil {
			u.Log.WithError(err).Error("failed to publish address created event")
			return nil, errkit.AddFuncName(err)
		}
		u.Log.Info("Published address created event")
	} else {
		u.Log.Info("Kafka producer is disabled, skipping address created event")
	}

	return converter.AddressToResponse(address), nil
}

func (u *AddressUseCaseImpl) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		u.Log.WithError(err).Error("failed to validate request body")
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, req.ContactId, req.UserId); err != nil {
		u.Log.WithError(err).Error("failed to find contact")
		return nil, errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIdAndContactId(tx, address, req.ID, contact.ID); err != nil {
		u.Log.WithError(err).Error("failed to find address")
		return nil, errkit.AddFuncName(err)
	}

	address.Street = req.Street
	address.City = req.City
	address.Province = req.Province
	address.PostalCode = req.PostalCode
	address.Country = req.Country

	if err := u.AddressRepository.Update(tx, address); err != nil {
		u.Log.WithError(err).Error("failed to update address")
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("failed to commit transaction")
		return nil, errkit.AddFuncName(err)
	}

	event := converter.AddressToEvent(address)
	if err := u.AddressProducer.Send(event); err != nil {
		u.Log.WithError(err).Error("failed to publish address updated event")
		return nil, errkit.AddFuncName(err)
	}

	return converter.AddressToResponse(address), nil
}

func (u *AddressUseCaseImpl) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, req.ContactId, req.UserId); err != nil {
		u.Log.WithError(err).Error("failed to find contact")
		return nil, errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIdAndContactId(tx, address, req.ID, req.ContactId); err != nil {
		u.Log.WithError(err).Error("failed to find address")
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("failed to commit transaction")
		return nil, errkit.AddFuncName(err)
	}

	return converter.AddressToResponse(address), nil
}

func (u *AddressUseCaseImpl) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, req.ContactId, req.UserId); err != nil {
		u.Log.WithError(err).Error("failed to find contact")
		return errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIdAndContactId(tx, address, req.ID, req.ContactId); err != nil {
		u.Log.WithError(err).Error("failed to find address")
		return errkit.AddFuncName(err)
	}

	if err := u.AddressRepository.Delete(tx, address); err != nil {
		u.Log.WithError(err).Error("failed to delete address")
		return errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("failed to commit transaction")
		return errkit.AddFuncName(err)
	}

	return nil
}

func (u *AddressUseCaseImpl) List(ctx context.Context, req *model.ListAddressRequest) ([]model.AddressResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, req.ContactId, req.UserId); err != nil {
		u.Log.WithError(err).Error("failed to find contact")
		return nil, errkit.AddFuncName(err)
	}

	addresses, err := u.AddressRepository.FindAllByContactId(tx, contact.ID)
	if err != nil {
		u.Log.WithError(err).Error("failed to find addresses")
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("failed to commit transaction")
		return nil, errkit.AddFuncName(err)
	}

	res := make([]model.AddressResponse, len(addresses))
	for i, address := range addresses {
		res[i] = *converter.AddressToResponse(&address)
	}

	return res, nil
}
