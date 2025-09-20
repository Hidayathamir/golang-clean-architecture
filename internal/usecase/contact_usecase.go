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
	Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error)
	Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error)
	Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error)
	Delete(ctx context.Context, req *model.DeleteContactRequest) error
	Search(ctx context.Context, req *model.SearchContactRequest) ([]model.ContactResponse, int64, error)
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

func (u *ContactUseCaseImpl) Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := &entity.Contact{
		ID:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		UserId:    req.UserId,
	}

	if err := u.ContactRepository.Create(tx, contact); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.ContactToEvent(contact)
	if err := u.ContactProducer.Send(event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.ContactToResponse(contact), nil
}

func (u *ContactUseCaseImpl) Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, req.ID, req.UserId); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact.FirstName = req.FirstName
	contact.LastName = req.LastName
	contact.Email = req.Email
	contact.Phone = req.Phone

	if err := u.ContactRepository.Update(tx, contact); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.ContactToEvent(contact)
	if err := u.ContactProducer.Send(event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.ContactToResponse(contact), nil
}

func (u *ContactUseCaseImpl) Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, req.ID, req.UserId); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.ContactToResponse(contact), nil
}

func (u *ContactUseCaseImpl) Delete(ctx context.Context, req *model.DeleteContactRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(tx, contact, req.ID, req.UserId); err != nil {
		return errkit.AddFuncName(err)
	}

	if err := u.ContactRepository.Delete(tx, contact); err != nil {
		return errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}

func (u *ContactUseCaseImpl) Search(ctx context.Context, req *model.SearchContactRequest) ([]model.ContactResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, 0, errkit.AddFuncName(err)
	}

	contacts, total, err := u.ContactRepository.Search(tx, req)
	if err != nil {
		return nil, 0, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, errkit.AddFuncName(err)
	}

	res := make([]model.ContactResponse, len(contacts))
	for i, contact := range contacts {
		res[i] = *converter.ContactToResponse(&contact)
	}

	return res, total, nil
}
