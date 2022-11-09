package kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"gitlab.com/dipper-iot/shared/logger"
	"strings"
	"time"
)

// Config kafka config
type Config struct {
	Brokers    []string `mapstructure:"brokers"`
	GroupID    string   `mapstructure:"groupID"`
	InitTopics bool     `mapstructure:"initTopics"`
}

// TopicConfig kafka topic config
type TopicConfig struct {
	TopicName         string `mapstructure:"topicName"`
	Partitions        int    `mapstructure:"partitions"`
	ReplicationFactor int    `mapstructure:"replicationFactor"`
}

type Options struct {
	minBytes               int
	maxBytes               int
	queueCapacity          int
	heartbeatInterval      time.Duration
	commitInterval         time.Duration
	partitionWatchInterval time.Duration
	maxAttempts            int
	dialTimeout            time.Duration
	maxWait                time.Duration
	writerReadTimeout      time.Duration
	writerWriteTimeout     time.Duration
	writerRequiredAcks     kafka.RequiredAcks
	writerMaxAttempts      int
	groupID                string
	errLogger              kafka.Logger
	balancer               kafka.Balancer
	compression            kafka.Compression
	async                  bool
	prefix                 string
}

func NewOptions() *Options {
	return &Options{
		minBytes:               10e3, // 10KB
		maxBytes:               10e6, // 10MB
		queueCapacity:          100,
		heartbeatInterval:      3 * time.Second,
		commitInterval:         0,
		partitionWatchInterval: 5 * time.Second,
		maxAttempts:            3,
		dialTimeout:            3 * time.Minute,
		maxWait:                1 * time.Second,
		writerReadTimeout:      10 * time.Second,
		writerWriteTimeout:     10 * time.Second,
		writerRequiredAcks:     -1,
		writerMaxAttempts:      3,
		errLogger:              kafka.LoggerFunc(logger.Errorf),
		balancer:               &kafka.LeastBytes{},
		compression:            kafka.Snappy,
		async:                  false,
		prefix:                 "",
		groupID:                "",
	}
}

func (o Options) Clone() Options {
	return o
}

type Option func(options *Options)

func GroupId(value string) Option {
	return func(options *Options) {
		options.groupID = value
	}
}

func PrefixConfig(prefix string) Option {
	if len(prefix) > 0 {
		prefix = fmt.Sprintf("%s_", strings.ToUpper(prefix))
	}
	return func(options *Options) {
		options.prefix = prefix
	}
}

func Balancer(value bool) Option {
	return func(options *Options) {
		options.async = value
	}
}

func Async(value kafka.Balancer) Option {
	return func(options *Options) {
		options.balancer = value
	}
}

func Compression(value kafka.Compression) Option {
	return func(options *Options) {
		options.compression = value
	}
}

func MinBytes(value int) Option {
	return func(options *Options) {
		options.minBytes = value
	}
}

func ErrorLogger(errLogger kafka.Logger) Option {
	return func(options *Options) {
		options.errLogger = errLogger
	}
}

func MaxBytes(value int) Option {
	return func(options *Options) {
		options.maxBytes = value
	}
}

func QueueCapacity(value int) Option {
	return func(options *Options) {
		options.queueCapacity = value
	}
}

func HeartbeatInterval(value time.Duration) Option {
	return func(options *Options) {
		options.heartbeatInterval = value
	}
}

func CommitInterval(value time.Duration) Option {
	return func(options *Options) {
		options.commitInterval = value
	}
}

func PartitionWatchInterval(value time.Duration) Option {
	return func(options *Options) {
		options.partitionWatchInterval = value
	}
}

func MaxAttempts(value int) Option {
	return func(options *Options) {
		options.maxAttempts = value
	}
}

func DialTimeout(value time.Duration) Option {
	return func(options *Options) {
		options.dialTimeout = value
	}
}

func MaxWait(value time.Duration) Option {
	return func(options *Options) {
		options.maxWait = value
	}
}

func WriterReadTimeout(value time.Duration) Option {
	return func(options *Options) {
		options.writerReadTimeout = value
	}
}
func WriterWriteTimeout(value time.Duration) Option {
	return func(options *Options) {
		options.writerWriteTimeout = value
	}
}

func WriterRequiredAcks(value kafka.RequiredAcks) Option {
	return func(options *Options) {
		options.writerRequiredAcks = value
	}
}

func GroupID(value string) Option {
	return func(options *Options) {
		options.groupID = value
	}
}
