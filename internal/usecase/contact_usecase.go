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

type ContactUseCase interface {
	Create(ctx context.Context, request *model.CreateContactRequest) (*model.ContactResponse, error)
	Update(ctx context.Context, request *model.UpdateContactRequest) (*model.ContactResponse, error)
	Get(ctx context.Context, request *model.GetContactRequest) (*model.ContactResponse, error)
	Delete(ctx context.Context, request *model.DeleteContactRequest) error
	Search(ctx context.Context, request *model.SearchContactRequest) ([]model.ContactResponse, int64, error)
}

var _ ContactUseCase = &ContactUseCaseImpl{}

type ContactUseCaseImpl struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate

	// repository
	ContactRepository repository.ContactRepository

	// producer
	ContactProducer messaging.ContactProducer
}

func NewContactUseCase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,

	// repository
	contactRepository repository.ContactRepository,

	// producer
	contactProducer messaging.ContactProducer,
) *ContactUseCaseImpl {
	return &ContactUseCaseImpl{
		DB:       db,
		Log:      logger,
		Validate: validate,

		// repository
		ContactRepository: contactRepository,

		// producer
		ContactProducer: contactProducer,
	}
}

func (u *ContactUseCaseImpl) Create(ctx context.Context, request *model.CreateContactRequest) (*model.ContactResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := &entity.Contact{
		ID:        uuid.New().String(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Phone:     request.Phone,
		UserId:    request.UserId,
	}

	if err := u.ContactRepository.Create(tx, contact); err != nil {
		u.Log.WithError(err).Error("error creating contact")
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error creating contact")
		return nil, errkit.AddFuncName(err)
	}

	event := converter.ContactToEvent(contact)
	if err := u.ContactProducer.Send(event); err != nil {
		u.Log.WithError(err).Error("error publishing contact created event")
		return nil, errkit.AddFuncName(err)
	}

	return converter.ContactToResponse(contact), nil
}

func (u *ContactUseCaseImpl) Update(ctx context.Context, request *model.UpdateContactRequest) (*model.ContactResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, request.ID, request.UserId); err != nil {
		u.Log.WithError(err).Error("error getting contact")
		return nil, errkit.AddFuncName(err)
	}

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact.FirstName = request.FirstName
	contact.LastName = request.LastName
	contact.Email = request.Email
	contact.Phone = request.Phone

	if err := u.ContactRepository.Update(tx, contact); err != nil {
		u.Log.WithError(err).Error("error updating contact")
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error updating contact")
		return nil, errkit.AddFuncName(err)
	}

	event := converter.ContactToEvent(contact)
	if err := u.ContactProducer.Send(event); err != nil {
		u.Log.WithError(err).Error("error publishing contact updated event")
		return nil, errkit.AddFuncName(err)
	}

	return converter.ContactToResponse(contact), nil
}

func (u *ContactUseCaseImpl) Get(ctx context.Context, request *model.GetContactRequest) (*model.ContactResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, request.ID, request.UserId); err != nil {
		u.Log.WithError(err).Error("error getting contact")
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error getting contact")
		return nil, errkit.AddFuncName(err)
	}

	return converter.ContactToResponse(contact), nil
}

func (u *ContactUseCaseImpl) Delete(ctx context.Context, request *model.DeleteContactRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, request.ID, request.UserId); err != nil {
		u.Log.WithError(err).Error("error getting contact")
		return errkit.AddFuncName(err)
	}

	if err := u.ContactRepository.Delete(tx, contact); err != nil {
		u.Log.WithError(err).Error("error deleting contact")
		return errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error deleting contact")
		return errkit.AddFuncName(err)
	}

	return nil
}

func (u *ContactUseCaseImpl) Search(ctx context.Context, request *model.SearchContactRequest) ([]model.ContactResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		err = errkit.BadRequest(err)
		return nil, 0, errkit.AddFuncName(err)
	}

	contacts, total, err := u.ContactRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Error("error getting contacts")
		return nil, 0, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error getting contacts")
		return nil, 0, errkit.AddFuncName(err)
	}

	responses := make([]model.ContactResponse, len(contacts))
	for i, contact := range contacts {
		responses[i] = *converter.ContactToResponse(&contact)
	}

	return responses, total, nil
}
