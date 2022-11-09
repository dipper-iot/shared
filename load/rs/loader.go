package rs

import (
	"fmt"
	base "github.com/go-redis/redis/v8"
	"gitlab.com/dipper-iot/shared/cache"
	"gitlab.com/dipper-iot/shared/cache/redis"
	"gitlab.com/dipper-iot/shared/cache/redis_sync"
	"gitlab.com/dipper-iot/shared/cli"
	"gitlab.com/dipper-iot/shared/service"
	"os"
	"strings"
)

const (
	Master string = ""
	Slave  string = "SLAVE"
)

type Redis struct {
	Config *RedisConfig
	client *base.Client
	prefix string
}

func (r *Redis) Flags() []cli.Flag {
	return nil
}

func (r *Redis) Priority() int {
	return 1
}

func (r *Redis) Start(o *service.Options, c *cli.Context) error {
	var (
		err error
	)

	address := os.Getenv(fmt.Sprintf("%sREDIS_ADDRESS", r.prefix))
	password := os.Getenv(fmt.Sprintf("%sREDIS_PASS", r.prefix))

	r.Config = &RedisConfig{
		Address:  address,
		Password: password,
	}

	r.client, err = NewRedisClient(r.Config.Address, r.Config.Password)
	return err
}

func NewRedis(prefix string) *Redis {
	if len(prefix) > 0 {
		prefix = fmt.Sprintf("%s_", strings.ToUpper(prefix))
	}
	return &Redis{
		Config: &RedisConfig{},
		prefix: prefix,
	}
}

func (r *Redis) Name() string {
	return "redis"
}

func (r *Redis) Stop() error {
	return r.client.Close()
}

func (r *Redis) Client() *base.Client {
	return r.client
}

func (r *Redis) Cache() cache.Cache {
	return redis.NewCache(r.client)
}

func (r *Redis) CacheSync() cache.Cache {
	return redis_sync.NewCache(r.client)
}
