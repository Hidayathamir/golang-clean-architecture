package contact_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestContactUsecaseImpl_Update_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
		SlackClient:       SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.UpdateContactRequest{
		UserId:    "userid1",
		ID:        uuid.NewString(),
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	ContactRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected = &model.ContactResponse{
		ID:        "",
		FirstName: "firstname1",
		LastName:  "",
		Email:     "hidayat@gmail.com",
		Phone:     "",
		CreatedAt: 0,
		UpdatedAt: 0,
		Addresses: nil,
	}

	assert.Equal(t, expected, res)
	assert.Nil(t, err)
}

func TestContactUsecaseImpl_Update_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
		SlackClient:       SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.UpdateContactRequest{
		UserId:    "",
		ID:        uuid.NewString(),
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	ContactRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestContactUsecaseImpl_Update_Fail_FindByIdAndUserId(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
		SlackClient:       SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.UpdateContactRequest{
		UserId:    "userid1",
		ID:        uuid.NewString(),
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return assert.AnError
	}

	ContactRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestContactUsecaseImpl_Update_Fail_Update(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
		SlackClient:       SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.UpdateContactRequest{
		UserId:    "userid1",
		ID:        uuid.NewString(),
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	ContactRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return assert.AnError
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestContactUsecaseImpl_Update_Fail_IsConnected(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
		SlackClient:       SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.UpdateContactRequest{
		UserId:    "userid1",
		ID:        uuid.NewString(),
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	ContactRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return false, assert.AnError
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestContactUsecaseImpl_Update_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ContactRepository := &mock.ContactRepositoryMock{}
	ContactProducer := &mock.ContactProducerMock{}
	SlackClient := &mock.SlackClientMock{}
	u := &contact.ContactUsecaseImpl{
		DB:                gormDB,
		Validate:          validator.New(),
		ContactRepository: ContactRepository,
		ContactProducer:   ContactProducer,
		SlackClient:       SlackClient,
	}

	// ------------------------------------------------------- //

	req := &model.UpdateContactRequest{
		UserId:    "userid1",
		ID:        uuid.NewString(),
		FirstName: "firstname1",
		Email:     "hidayat@gmail.com",
	}

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	ContactRepository.UpdateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Contact) error {
		return nil
	}

	SlackClient.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	ContactProducer.SendFunc = func(ctx context.Context, event *model.ContactEvent) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	res, err := u.Update(context.Background(), req)

	// ------------------------------------------------------- //

	var expected *model.ContactResponse

	assert.Equal(t, expected, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
