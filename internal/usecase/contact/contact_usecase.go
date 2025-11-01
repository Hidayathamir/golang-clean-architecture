package contact

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/ContactUsecase.go -pkg=mock . ContactUsecase

type ContactUsecase interface {
	Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error)
	Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error)
	Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error)
	Delete(ctx context.Context, req *model.DeleteContactRequest) error
	Search(ctx context.Context, req *model.SearchContactRequest) (model.ContactResponseList, int64, error)
}

var _ ContactUsecase = &ContactUsecaseImpl{}

type ContactUsecaseImpl struct {
	Config   *viper.Viper
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate

	// repository
	ContactRepository repository.ContactRepository

	// producer
	ContactProducer messaging.ContactProducer

	// client
	SlackClient rest.SlackClient
}

func NewContactUsecase(
	cfg *viper.Viper, log *logrus.Logger, db *gorm.DB, validate *validator.Validate,

	// repository
	contactRepository repository.ContactRepository,

	// producer
	contactProducer messaging.ContactProducer,

	// client
	SlackClient rest.SlackClient,
) *ContactUsecaseImpl {
	return &ContactUsecaseImpl{
		Config:   cfg,
		Log:      log,
		DB:       db,
		Validate: validate,

		// repository
		ContactRepository: contactRepository,

		// producer
		ContactProducer: contactProducer,

		// client
		SlackClient: SlackClient,
	}
}
