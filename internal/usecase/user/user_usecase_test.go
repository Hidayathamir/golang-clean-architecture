package user_test

import (
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/storage"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewUserUsecase(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	var DB = gormDB
	var Config = viper.New()

	var UserRepository repository.UserRepository = &mock.UserRepositoryMock{}
	var FollowRepository repository.FollowRepository = &mock.FollowRepositoryMock{}

	var UserProducer messaging.UserProducer = &mock.UserProducerMock{}
	var NotifProducer messaging.NotifProducer = &mock.NotifProducerMock{}

	var S3Client storage.S3Client = &mock.S3ClientMock{}

	u := user.NewUserUsecase(Config, DB, UserRepository, FollowRepository, UserProducer, NotifProducer, S3Client)

	require.NotEmpty(t, u)
}
