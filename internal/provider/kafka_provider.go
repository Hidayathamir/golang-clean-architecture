package provider

import (
	"strings"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kotel"
)

func NewKafkaClientProducer(cfg *config.Config) *kgo.Client {
	brokers := strings.Split(cfg.GetKafkaBootstrapServers(), ",")
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.RequestRetries(20),
	}

	tracer := kotel.NewTracer()
	opts = append(opts, kgo.WithHooks(tracer))

	client, err := kgo.NewClient(opts...)
	x.PanicIfErr(err)

	return client
}

func NewKafkaClientConsumer(cfg *config.Config, consumerGroup string, topic string) *kgo.Client {
	brokers := strings.Split(cfg.GetKafkaBootstrapServers(), ",")
	const _1MB = 1024 * 1024
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(consumerGroup),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.DisableAutoCommit(),
		kgo.FetchMinBytes(_1MB),           // will wait atleast 1 MB before processing
		kgo.FetchMaxWait(5 * time.Second), // but with interval 5 second
	}

	tracer := kotel.NewTracer()
	opts = append(opts, kgo.WithHooks(tracer))

	client, err := kgo.NewClient(opts...)
	x.PanicIfErr(err)

	return client
}
