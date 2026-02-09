package provider

import (
	"strings"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
)

func NewKafkaConsumerGroup(config *config.Config, groupID string) sarama.ConsumerGroup {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	offsetReset := config.GetString(configkey.KafkaAutoOffsetReset)
	if offsetReset == "earliest" {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	brokers := strings.Split(config.GetString(configkey.KafkaBootstrapServers), ",")

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		x.Logger.Panicf("Failed to create consumer group: %v", err)
	}
	return consumerGroup
}

func NewKafkaProducer(cfg *config.Config) sarama.SyncProducer {
	if !cfg.GetBool(configkey.KafkaProducerEnabled) {
		x.Logger.Info("Kafka producer is disabled")
		return nil
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3

	brokers := strings.Split(cfg.GetString(configkey.KafkaBootstrapServers), ",")

	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		x.Logger.Panicf("Failed to create producer: %v", err)
	}

	producer = otelsarama.WrapSyncProducer(saramaConfig, producer)

	return producer
}
