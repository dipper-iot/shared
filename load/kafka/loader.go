package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gitlab.com/dipper-iot/shared/cli"
	"gitlab.com/dipper-iot/shared/service"
	"os"
	"strings"
)

type Kafka struct {
	prefix        string
	Brokers       []string
	producers     []Producer
	consumerGroup *consumerGroup
	options       *Options
}

func NewKafka(opts ...Option) *Kafka {
	options := NewOptions()

	for _, opt := range opts {
		opt(options)
	}

	return &Kafka{
		prefix:    options.prefix,
		Brokers:   []string{"localhost:9200"},
		producers: []Producer{},
		options:   options,
	}
}

func (k Kafka) Name() string {
	return "kafka"
}

func (k Kafka) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (k Kafka) Priority() int {
	return 1
}

func (k *Kafka) Start(o *service.Options, c *cli.Context) error {

	broker := os.Getenv(fmt.Sprintf("%sKAFKA_BROKERS", k.prefix))
	k.Brokers = strings.Split(broker, ",")
	k.consumerGroup = NewConsumerGroup(k.Brokers)

	return nil
}

func (k Kafka) Stop() error {
	for _, p := range k.producers {
		p.Close()
	}

	return nil
}

func (k Kafka) ConsumeGroupTopic(ctx context.Context, groupTopics []string, poolSize int, worker Worker, opts ...Option) {
	options := k.getOptions(opts)

	k.consumerGroup.ConsumeTopic(ctx, groupTopics, poolSize, worker, options)
}

func (k Kafka) ConsumeTopic(topic string, opts ...Option) *kafka.Reader {
	options := k.getOptions(opts)

	return NewKafkaReader(k.Brokers, topic, options)
}

func (k *Kafka) Producer(opts ...Option) Producer {
	options := k.getOptions(opts)

	producer := NewProducer(k.Brokers, options)
	k.producers = append(k.producers, producer)
	return producer
}

func (k Kafka) getOptions(opts []Option) *Options {
	options := k.options.Clone()

	for _, opt := range opts {
		opt(&options)
	}
	return &options
}
