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
	"gorm.io/gorm"
)

func TestUserUsecaseImpl_Create_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
	}

	// ------------------------------------------------------- //

	req := &model.RegisterUserRequest{
		ID:       "id1",
		Password: "pw1",
		Name:     "name1",
	}

	UserRepository.CountByIdFunc = func(ctx context.Context, db *gorm.DB, id string) (int64, error) {
		return 0, nil
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse = &model.UserResponse{
		ID:        "id1",
		Name:      "name1",
		Token:     "",
		CreatedAt: 0,
		UpdatedAt: 0,
	}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestUserUsecaseImpl_Create_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
	}

	// ------------------------------------------------------- //

	req := &model.RegisterUserRequest{
		ID:       "",
		Password: "pw1",
		Name:     "name1",
	}

	UserRepository.CountByIdFunc = func(ctx context.Context, db *gorm.DB, id string) (int64, error) {
		return 0, nil
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Create_Fail_CountById(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
	}

	// ------------------------------------------------------- //

	req := &model.RegisterUserRequest{
		ID:       "id1",
		Password: "pw1",
		Name:     "name1",
	}

	UserRepository.CountByIdFunc = func(ctx context.Context, db *gorm.DB, id string) (int64, error) {
		return 0, assert.AnError
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Create_Fail_UserAlreadyExists(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
	}

	// ------------------------------------------------------- //

	req := &model.RegisterUserRequest{
		ID:       "id1",
		Password: "pw1",
		Name:     "name1",
	}

	UserRepository.CountByIdFunc = func(ctx context.Context, db *gorm.DB, id string) (int64, error) {
		return 1, nil
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "already exists")
}

func TestUserUsecaseImpl_Create_Fail_Create(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
	}

	// ------------------------------------------------------- //

	req := &model.RegisterUserRequest{
		ID:       "id1",
		Password: "pw1",
		Name:     "name1",
	}

	UserRepository.CountByIdFunc = func(ctx context.Context, db *gorm.DB, id string) (int64, error) {
		return 0, nil
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return assert.AnError
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Create_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
	}

	// ------------------------------------------------------- //

	req := &model.RegisterUserRequest{
		ID:       "id1",
		Password: "pw1",
		Name:     "name1",
	}

	UserRepository.CountByIdFunc = func(ctx context.Context, db *gorm.DB, id string) (int64, error) {
		return 0, nil
	}

	UserRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Create(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
