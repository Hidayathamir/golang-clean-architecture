package contact_test

import (
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewContactUsecase(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	var DB = gormDB
	var Log = logrus.New()
	var Validate = validator.New()

	var ContactRepository repository.ContactRepository = &mock.ContactRepositoryMock{}

	var ContactProducer messaging.ContactProducer = &mock.ContactProducerMock{}

	var SlackClient rest.SlackClient = &mock.SlackClientMock{}

	u := contact.NewContactUsecase(
		DB, Log, Validate,

		ContactRepository,

		ContactProducer,

		SlackClient,
	)

	assert.NotEmpty(t, u)
}
