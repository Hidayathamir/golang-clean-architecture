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

func TestUserUsecaseImpl_Logout_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	S3Client := &mock.S3ClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		S3Client:       S3Client,
	}

	// ------------------------------------------------------- //

	req := &model.LogoutUserRequest{
		ID: "id1",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	S3Client.DeleteObjectFunc = func(ctx context.Context, bucket, key string) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Logout(context.Background(), req)

	// ------------------------------------------------------- //

	var expected = true

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestUserUsecaseImpl_Logout_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	S3Client := &mock.S3ClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		S3Client:       S3Client,
	}

	// ------------------------------------------------------- //

	req := &model.LogoutUserRequest{
		ID: "",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	S3Client.DeleteObjectFunc = func(ctx context.Context, bucket, key string) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Logout(context.Background(), req)

	// ------------------------------------------------------- //

	var expected bool

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Logout_Fail_FindByID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	S3Client := &mock.S3ClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		S3Client:       S3Client,
	}

	// ------------------------------------------------------- //

	req := &model.LogoutUserRequest{
		ID: "id1",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return assert.AnError
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	S3Client.DeleteObjectFunc = func(ctx context.Context, bucket, key string) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Logout(context.Background(), req)

	// ------------------------------------------------------- //

	var expected bool

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Logout_Fail_Update(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	S3Client := &mock.S3ClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		S3Client:       S3Client,
	}

	// ------------------------------------------------------- //

	req := &model.LogoutUserRequest{
		ID: "id1",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return assert.AnError
	}

	S3Client.DeleteObjectFunc = func(ctx context.Context, bucket, key string) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Logout(context.Background(), req)

	// ------------------------------------------------------- //

	var expected bool

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Logout_Fail_DeleteObject(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	S3Client := &mock.S3ClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		S3Client:       S3Client,
	}

	// ------------------------------------------------------- //

	req := &model.LogoutUserRequest{
		ID: "id1",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	S3Client.DeleteObjectFunc = func(ctx context.Context, bucket, key string) (bool, error) {
		return false, assert.AnError
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Logout(context.Background(), req)

	// ------------------------------------------------------- //

	var expected bool

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Logout_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	UserRepository := &mock.UserRepositoryMock{}
	UserProducer := &mock.UserProducerMock{}
	S3Client := &mock.S3ClientMock{}
	u := &user.UserUsecaseImpl{
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: UserRepository,
		UserProducer:   UserProducer,
		S3Client:       S3Client,
	}

	// ------------------------------------------------------- //

	req := &model.LogoutUserRequest{
		ID: "id1",
	}

	UserRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	UserRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User) error {
		return nil
	}

	S3Client.DeleteObjectFunc = func(ctx context.Context, bucket, key string) (bool, error) {
		return true, nil
	}

	UserProducer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Logout(context.Background(), req)

	// ------------------------------------------------------- //

	var expected bool

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
