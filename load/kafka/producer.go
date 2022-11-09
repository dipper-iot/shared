package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	brokers []string
	w       *kafka.Writer
}

// NewProducer create new kafka producer
func NewProducer(brokers []string, opts *Options) *producer {
	return &producer{brokers: brokers, w: NewWriter(brokers, opts)}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *producer) Close() error {
	return p.w.Close()
}
