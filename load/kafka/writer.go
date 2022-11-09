package kafka

import (
	"github.com/segmentio/kafka-go"
)

// NewWriter create new configured kafka writer
func NewWriter(brokers []string, options *Options) *kafka.Writer {

	w := &kafka.Writer{
		Addr: kafka.TCP(brokers...),

		Balancer:     options.balancer,
		RequiredAcks: options.writerRequiredAcks,
		MaxAttempts:  options.writerMaxAttempts,
		ErrorLogger:  options.errLogger,
		Compression:  options.compression,
		ReadTimeout:  options.writerReadTimeout,
		WriteTimeout: options.writerWriteTimeout,
		Async:        options.async,
	}
	return w
}
