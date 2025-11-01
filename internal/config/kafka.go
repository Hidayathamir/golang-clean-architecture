package config

import (
	"strings"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewKafkaConsumerGroup(config *viper.Viper, log *logrus.Logger) sarama.ConsumerGroup {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	offsetReset := config.GetString(configkey.KafkaAutoOffsetReset)
	if offsetReset == "earliest" {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	brokers := strings.Split(config.GetString(configkey.KafkaBootstrapServers), ",")
	groupID := config.GetString(configkey.KafkaGroupId)

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		log.Panicf("Failed to create consumer group: %v", err)
	}
	return consumerGroup
}

func NewKafkaProducer(viperConfig *viper.Viper, log *logrus.Logger) sarama.SyncProducer {
	if !viperConfig.GetBool(configkey.KafkaProducerEnabled) {
		log.Info("Kafka producer is disabled")
		return nil
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3

	brokers := strings.Split(viperConfig.GetString(configkey.KafkaBootstrapServers), ",")

	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		log.Panicf("Failed to create producer: %v", err)
	}

	producer = otelsarama.WrapSyncProducer(saramaConfig, producer)

	return producer
}
