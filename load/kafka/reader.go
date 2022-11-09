package kafka

import (
	"github.com/segmentio/kafka-go"
)

// NewKafkaReader create new configured kafka reader
func NewKafkaReader(kafkaURL []string, topic string, options *Options) *kafka.Reader {

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                options.groupID,
		Topic:                  topic,
		MinBytes:               options.minBytes,
		MaxBytes:               options.maxBytes,
		QueueCapacity:          options.queueCapacity,
		HeartbeatInterval:      options.heartbeatInterval,
		CommitInterval:         options.commitInterval,
		PartitionWatchInterval: options.partitionWatchInterval,
		ErrorLogger:            options.errLogger,
		MaxAttempts:            options.maxAttempts,
		MaxWait:                options.maxWait,
		Dialer: &kafka.Dialer{
			Timeout: options.dialTimeout,
		},
	})

}
