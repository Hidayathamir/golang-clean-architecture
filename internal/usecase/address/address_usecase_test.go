package address_test

import (
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewAddressUsecase(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	var DB *gorm.DB = gormDB
	var Log *logrus.Logger = logrus.New()
	var Validate *validator.Validate = validator.New()

	var AddressRepository repository.AddressRepository = &mock.AddressRepositoryMock{}
	var ContactRepository repository.ContactRepository = &mock.ContactRepositoryMock{}

	var AddressProducer messaging.AddressProducer = &mock.AddressProducerMock{}

	var PaymentClient rest.PaymentClient = &mock.PaymentClientMock{}

	u := address.NewAddressUsecase(
		DB, Log, Validate,

		ContactRepository,
		AddressRepository,

		AddressProducer,

		PaymentClient,
	)

	assert.NotEmpty(t, u)
}
