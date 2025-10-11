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

func TestUserUsecaseImpl_Update_Success(t *testing.T) {
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

	req := &model.UpdateUserRequest{
		ID:       "id1",
		Password: "password1",
		Name:     "name1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected = &model.UserResponse{Name: "name1"}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestUserUsecaseImpl_Update_Fail_ValidateStruct(t *testing.T) {
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

	req := &model.UpdateUserRequest{
		ID: "",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Update_Fail_FindById(t *testing.T) {
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

	req := &model.UpdateUserRequest{
		ID: "id1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return assert.AnError
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Update_Fail_Update(t *testing.T) {
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

	req := &model.UpdateUserRequest{
		ID: "id1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return assert.AnError
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Update_Fail_Send(t *testing.T) {
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

	req := &model.UpdateUserRequest{
		ID: "id1",
	}

	UserRepository.FindByIdFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.UserResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
