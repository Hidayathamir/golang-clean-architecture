package user_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_Login_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		SlackClient:    SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Login(context.Background(), req)

	// ------------------------------------------------------- //

	res.Token = ""
	var expected *model.UserResponse = &model.UserResponse{}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestUserUsecaseImpl_Login_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		SlackClient:    SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.LoginUserRequest{
		ID:       "",
		Password: "password1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Login(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Login_Fail_FindById(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		SlackClient:    SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return assert.AnError
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Login(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Login_Fail_CompareHashAndPassword(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		SlackClient:    SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return assert.AnError
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Login(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
}

func TestUserUsecaseImpl_Login_Fail_Update(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		SlackClient:    SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return assert.AnError
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Login(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Login_Fail_IsConnected(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		SlackClient:    SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return false, assert.AnError
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Login(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Login_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		SlackClient:    SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Login(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
