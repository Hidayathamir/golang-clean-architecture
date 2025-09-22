package address

import (
	"context"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/gateway/rest"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
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

	// client
	PaymentClient rest.PaymentClient
}

func NewAddressUseCase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,

	// repository
	contactRepository repository.ContactRepository,
	addressRepository repository.AddressRepository,

	// producer
	addressProducer messaging.AddressProducer,

	// client
	paymentClient rest.PaymentClient,
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

		// client
		PaymentClient: paymentClient,
	}
}
