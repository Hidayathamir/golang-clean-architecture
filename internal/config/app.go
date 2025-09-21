package config

import (
	"golang-clean-architecture/internal/delivery/http"
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/delivery/http/route"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/gateway/messaging/messagingmwlogger"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/internal/repository/repositorymwlogger"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/internal/usecase/usecasemwlogger"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Producer sarama.SyncProducer
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	var userRepository repository.UserRepository
	userRepository = repository.NewUserRepository(config.Log)
	userRepository = repositorymwlogger.NewUserRepository(config.Log, userRepository)

	var contactRepository repository.ContactRepository
	contactRepository = repository.NewContactRepository(config.Log)
	contactRepository = repositorymwlogger.NewContactRepository(config.Log, contactRepository)

	var addressRepository repository.AddressRepository
	addressRepository = repository.NewAddressRepository(config.Log)
	addressRepository = repositorymwlogger.NewAddressRepository(config.Log, addressRepository)

	// setup producer
	var userProducer messaging.UserProducer
	userProducer = messaging.NewUserProducer(config.Producer, config.Log)
	userProducer = messagingmwlogger.NewUserProducer(config.Log, userProducer)

	var contactProducer messaging.ContactProducer
	contactProducer = messaging.NewContactProducer(config.Producer, config.Log)
	contactProducer = messagingmwlogger.NewContactProducer(config.Log, contactProducer)

	var addressProducer messaging.AddressProducer
	addressProducer = messaging.NewAddressProducer(config.Producer, config.Log)
	addressProducer = messagingmwlogger.NewAddressProducer(config.Log, addressProducer)

	// setup use cases
	var userUseCase usecase.UserUseCase
	userUseCase = usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, userProducer)
	userUseCase = usecasemwlogger.NewUserUseCase(config.Log, userUseCase)

	var contactUseCase usecase.ContactUseCase
	contactUseCase = usecase.NewContactUseCase(config.DB, config.Log, config.Validate, contactRepository, contactProducer)
	contactUseCase = usecasemwlogger.NewContactUseCase(config.Log, contactUseCase)

	var addressUseCase usecase.AddressUseCase
	addressUseCase = usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, contactRepository, addressRepository, addressProducer)
	addressUseCase = usecasemwlogger.NewAddressUseCase(config.Log, addressUseCase)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	contactController := http.NewContactController(contactUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	traceIDMiddleware := middleware.NewTraceID()

	routeConfig := route.RouteConfig{
		App:               config.App,
		UserController:    userController,
		ContactController: contactController,
		AddressController: addressController,
		AuthMiddleware:    authMiddleware,
		TraceIDMiddleware: traceIDMiddleware,
	}

	routeConfig.Setup()
}
