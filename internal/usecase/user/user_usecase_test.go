package user_test

import (
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewUserUsecase(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	var DB *gorm.DB = gormDB
	var Log *logrus.Logger = logrus.New()
	var Validate *validator.Validate = validator.New()

	var UserRepository repository.UserRepository = &mock.UserRepositoryMock{}

	var UserProducer messaging.UserProducer = &mock.UserProducerMock{}

	var S3Client rest.S3Client = &mock.S3ClientMock{}
	var SlackClient rest.SlackClient = &mock.SlackClientMock{}

	u := user.NewUserUsecase(DB, Log, Validate, UserRepository, UserProducer, S3Client, SlackClient)

	assert.NotEmpty(t, u)
}
