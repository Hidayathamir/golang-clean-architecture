package address_test

import (
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewAddressUsecase(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	var DB = gormDB
	var Config = viper.New()

	var AddressRepository repository.AddressRepository = &mock.AddressRepositoryMock{}
	var ContactRepository repository.ContactRepository = &mock.ContactRepositoryMock{}

	var AddressProducer messaging.AddressProducer = &mock.AddressProducerMock{}

	var PaymentClient rest.PaymentClient = &mock.PaymentClientMock{}

	u := address.NewAddressUsecase(
		Config, DB,

		ContactRepository,
		AddressRepository,

		AddressProducer,

		PaymentClient,
	)

	assert.NotEmpty(t, u)
}
