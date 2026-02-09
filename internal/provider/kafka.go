package provider

import (
	"strings"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
)

func NewKafkaConsumerGroup(config *config.Config, groupID string) sarama.ConsumerGroup {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	offsetReset := config.GetKafkaAutoOffsetReset()
	if offsetReset == "" {
		offsetReset = "earliest"
	}
	if offsetReset == "earliest" {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	brokers := strings.Split(config.GetKafkaBootstrapServers(), ",")
	if len(brokers) == 0 {
		x.Logger.Panicf("Kafka brokers are not configured for consumer group %s", groupID)
	}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		x.Logger.Panicf("Failed to create consumer group: %v", err)
	}
	return consumerGroup
}

func NewKafkaProducer(cfg *config.Config) sarama.SyncProducer {
	if !cfg.GetKafkaProducerEnabled() {
		x.Logger.Info("Kafka producer is disabled")
		return nil
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3

	brokers := strings.Split(cfg.GetKafkaBootstrapServers(), ",")

	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		x.Logger.Panicf("Failed to create producer: %v", err)
	}

	producer = otelsarama.WrapSyncProducer(saramaConfig, producer)

	return producer
}
