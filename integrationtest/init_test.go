package integrationtest

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/route"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"
)

var app *fiber.App

var db *gorm.DB

var viperConfig *viper.Viper

var validate *validator.Validate

// TestMain is the entry point for all tests in this package.
// It sets up global dependencies (logger, validator, Fiber app, DB, Kafka producer),
// starts a PostgreSQL container, runs database migrations, bootstraps the application,
// executes test, and finally terminates the container before exiting.
func TestMain(m *testing.M) {
	viperConfig = config.NewViper()
	l.SetupLogger(viperConfig)
	validate = config.NewValidator(viperConfig)
	app = config.NewFiber(viperConfig)
	postgresContainer := newPostgresContainer(viperConfig)
	viperConfig.Set(configkey.DatabaseMigrations, "../db/migrations")
	config.Migrate(viperConfig)
	db = config.NewDatabase(viperConfig)
	kafkaContainer := newKafkaContainer(viperConfig)
	producer := config.NewKafkaProducer(viperConfig)

	usecases := config.SetupUsecases(viperConfig, db, app, validate, producer)
	controllers := config.SetupControllers(viperConfig, usecases)
	middlewares := config.SetupMiddlewares(usecases)

	route.Setup(app, controllers, middlewares)

	code := m.Run()

	panicIfErr(testcontainers.TerminateContainer(postgresContainer))
	panicIfErr(testcontainers.TerminateContainer(kafkaContainer))

	os.Exit(code)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func newPostgresContainer(viperConfig *viper.Viper) (postgresContainer *postgres.PostgresContainer) {
	database := viperConfig.GetString(configkey.DatabaseName)
	username := viperConfig.GetString(configkey.DatabaseUsername)
	password := viperConfig.GetString(configkey.DatabasePassword)

	var err error
	postgresContainer, err = postgres.Run(context.Background(),
		"postgres:15",
		postgres.WithDatabase(database),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
	)
	panicIfErr(err)
	state, err := postgresContainer.State(context.Background())
	panicIfErr(err)
	if !state.Running {
		panic("postgres container not running")
	}

	host, err := postgresContainer.Host(context.Background())
	panicIfErr(err)
	mappedPort, err := postgresContainer.MappedPort(context.Background(), "5432/tcp")
	panicIfErr(err)

	port, err := strconv.Atoi(mappedPort.Port())
	panicIfErr(err)
	viperConfig.Set(configkey.DatabaseHost, host)
	viperConfig.Set(configkey.DatabasePort, port)

	return postgresContainer
}

func newKafkaContainer(viperConfig *viper.Viper) (kafkaContainer *kafka.KafkaContainer) {
	var err error
	kafkaContainer, err = kafka.Run(context.Background(),
		"confluentinc/cp-kafka:7.2.15",
	)
	panicIfErr(err)
	state, err := kafkaContainer.State(context.Background())
	panicIfErr(err)
	if !state.Running {
		panic("kafka container not running")
	}

	brokers, err := kafkaContainer.Brokers(context.Background())
	panicIfErr(err)
	viperConfig.Set(configkey.KafkaBootstrapServers, brokers[0])

	return kafkaContainer
}
