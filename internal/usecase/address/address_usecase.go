package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/AddressUsecase.go -pkg=mock . AddressUsecase

type AddressUsecase interface {
	Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error)
	Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error)
	Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error)
	Delete(ctx context.Context, req *model.DeleteAddressRequest) error
	List(ctx context.Context, req *model.ListAddressRequest) (model.AddressResponseList, error)
}

var _ AddressUsecase = &AddressUsecaseImpl{}

type AddressUsecaseImpl struct {
	Config *viper.Viper
	DB     *gorm.DB

	// repository
	AddressRepository repository.AddressRepository
	ContactRepository repository.ContactRepository

	// producer
	AddressProducer messaging.AddressProducer

	// client
	PaymentClient rest.PaymentClient
}

func NewAddressUsecase(
	cfg *viper.Viper, db *gorm.DB,

	// repository
	contactRepository repository.ContactRepository,
	addressRepository repository.AddressRepository,

	// producer
	addressProducer messaging.AddressProducer,

	// client
	paymentClient rest.PaymentClient,
) *AddressUsecaseImpl {
	return &AddressUsecaseImpl{
		Config: cfg,
		DB:     db,

		// repository
		ContactRepository: contactRepository,
		AddressRepository: addressRepository,

		// producer
		AddressProducer: addressProducer,

		// client
		PaymentClient: paymentClient,
	}
}
