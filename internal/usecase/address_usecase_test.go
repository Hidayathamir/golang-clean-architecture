package usecase_test

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/mock"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newFakeDB() (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, db, _ := sqlmock.New()
	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	return gormDB, db
}

func TestAddressUseCaseImpl_Get_Success(t *testing.T) {
	gormDB, db := newFakeDB()
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &usecase.AddressUseCaseImpl{
		DB:                gormDB,
		Log:               logrus.New(),
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	db.ExpectBegin()
	db.ExpectCommit()

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return nil
	}

	// ------------------------------------------------------- //

	actual, err := u.Get(context.Background(), &model.GetAddressRequest{})
	assert.Nil(t, err)

	// ------------------------------------------------------- //

	expected := &model.AddressResponse{}

	assert.Equal(t, expected, actual)
}

func TestAddressUseCaseImpl_Get_Fail1(t *testing.T) {
	gormDB, _ := newFakeDB()
	AddressRepository := &mock.AddressRepositoryMock{}
	ContactRepository := &mock.ContactRepositoryMock{}
	u := &usecase.AddressUseCaseImpl{
		DB:                gormDB,
		Log:               logrus.New(),
		AddressRepository: AddressRepository,
		ContactRepository: ContactRepository,
	}

	// ------------------------------------------------------- //

	ContactRepository.FindByIdAndUserIdFunc = func(ctx context.Context, db *gorm.DB, contact *entity.Contact, id, userId string) error {
		return nil
	}

	AddressRepository.FindByIdAndContactIdFunc = func(ctx context.Context, db *gorm.DB, address *entity.Address, id, contactId string) error {
		return assert.AnError
	}

	// ------------------------------------------------------- //

	actual, err := u.Get(context.Background(), &model.GetAddressRequest{})
	assert.NotNil(t, err)

	// ------------------------------------------------------- //

	var expected *model.AddressResponse

	assert.Equal(t, expected, actual)
}
