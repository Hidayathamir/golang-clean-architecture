package contact

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ContactUsecase interface {
	Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error)
	Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error)
	Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error)
	Delete(ctx context.Context, req *model.DeleteContactRequest) error
	Search(ctx context.Context, req *model.SearchContactRequest) ([]model.ContactResponse, int64, error)
}

var _ ContactUsecase = &ContactUsecaseImpl{}

type ContactUsecaseImpl struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate

	// repository
	ContactRepository repository.ContactRepository

	// producer
	ContactProducer messaging.ContactProducer

	// client
	SlackClient rest.SlackClient
}

func NewContactUsecase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,

	// repository
	contactRepository repository.ContactRepository,

	// producer
	contactProducer messaging.ContactProducer,

	// client
	SlackClient rest.SlackClient,
) *ContactUsecaseImpl {
	return &ContactUsecaseImpl{
		DB:       db,
		Log:      logger,
		Validate: validate,

		// repository
		ContactRepository: contactRepository,

		// producer
		ContactProducer: contactProducer,

		// client
		SlackClient: SlackClient,
	}
}
