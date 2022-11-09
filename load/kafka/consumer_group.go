package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"gitlab.com/dipper-iot/shared/logger"
	"sync"
)

// MessageProcessor processor methods must implement kafka.Worker func method interface
type MessageProcessor interface {
	ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)
}

// Worker kafka consumer worker fetch and process messages from reader
type Worker func(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)

type ConsumerGroup interface {
	ConsumeTopic(ctx context.Context, cancel context.CancelFunc, groupID, topic string, poolSize int, worker Worker)
	GetNewKafkaReader(kafkaURL []string, topic, groupID string) *kafka.Reader
	GetNewKafkaWriter(topic string) *kafka.Writer
}

type consumerGroup struct {
	Brokers []string
}

// NewConsumerGroup kafka consumer group constructor
func NewConsumerGroup(brokers []string) *consumerGroup {
	return &consumerGroup{
		Brokers: brokers,
	}
}

// GetNewKafkaReader create new kafka reader
func (c *consumerGroup) getNewKafkaReader(kafkaURL []string, groupTopics []string, options *Options) *kafka.Reader {

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                options.groupID,
		GroupTopics:            groupTopics,
		MinBytes:               options.minBytes,
		MaxBytes:               options.maxBytes,
		QueueCapacity:          options.queueCapacity,
		HeartbeatInterval:      options.heartbeatInterval,
		CommitInterval:         options.commitInterval,
		PartitionWatchInterval: options.partitionWatchInterval,
		MaxAttempts:            options.maxAttempts,
		MaxWait:                options.maxWait,
		Dialer: &kafka.Dialer{
			Timeout: options.dialTimeout,
		},
	})
}

// ConsumeTopic start consumer group with given worker and pool size
func (c *consumerGroup) ConsumeTopic(ctx context.Context, groupTopics []string, poolSize int, worker Worker, options *Options) {

	r := c.getNewKafkaReader(c.Brokers, groupTopics, options)

	defer func() {
		if err := r.Close(); err != nil {
			logger.Warnf("consumerGroup.r.Close: %v", err)
		}
	}()

	logger.Infof("Starting consumer groupID: %s, topic: %+v, pool size: %v", options.groupID, groupTopics, poolSize)

	wg := &sync.WaitGroup{}
	for i := 0; i <= poolSize; i++ {
		wg.Add(1)
		go worker(ctx, r, wg, i)
	}
	wg.Wait()
}
