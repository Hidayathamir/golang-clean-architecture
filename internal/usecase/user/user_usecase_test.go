package user_test

import (
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/cache"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/storage"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/stretchr/testify/require"
)

func TestNewUserUsecase(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	var DB = gormDB
	var Config = config.NewConfig()

	var UserRepository repository.UserRepository = &mock.UserRepositoryMock{}
	var UserStatRepository repository.UserStatRepository = &mock.UserStatRepositoryMock{}
	var FollowRepository repository.FollowRepository = &mock.FollowRepositoryMock{}

	var UserProducer messaging.UserProducer = &mock.UserProducerMock{}
	var NotifProducer messaging.NotifProducer = &mock.NotifProducerMock{}

	var S3Client storage.S3Client = &mock.S3ClientMock{}
	var UserCache cache.UserCache = &mock.UserCacheMock{}

	u := user.NewUserUsecase(Config, DB, UserRepository, UserStatRepository, FollowRepository, UserProducer, NotifProducer, S3Client, UserCache)

	require.NotEmpty(t, u)
}
