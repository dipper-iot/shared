package redis

import (
	"context"
	"encoding/json"
	base "github.com/go-redis/redis/v8"
	"gitlab.com/dipper-iot/shared/cache"
	"gitlab.com/dipper-iot/shared/load/rs"
	"gitlab.com/dipper-iot/shared/logger"
	"time"
)

type redisCache struct {
	client  *base.Client
	options *cache.OptionsCache
}

func NewCache(client *base.Client) *redisCache {
	return &redisCache{
		client:  client,
		options: cache.NewOptionsCache(),
	}
}

func NewCacheEnv() (*redisCache, error) {
	client, err := rs.NewRedisToEnv()
	if err != nil {
		return nil, err
	}
	return &redisCache{
		client:  client,
		options: cache.NewOptionsCache(),
	}, nil
}

func (r *redisCache) Init(options ...cache.OptionCache) error {
	r.options.Load(options)
	return nil
}

func (r *redisCache) Get(ctx context.Context, key string, data interface{}, options ...cache.OptionCache) (bool, error) {
	logger.Tracef("Redis cache get %s", key)
	jsonStr, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return false, nil
	}

	if len(jsonStr) > 0 {
		err = json.Unmarshal([]byte(jsonStr), data)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (r *redisCache) Set(ctx context.Context, key string, data interface{}, options ...cache.OptionCache) error {
	bData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	logger.Tracef("Redis cache set %s - %s", key, string(bData))

	timeCache, success := cache.TimeCacheFromContext(ctx)
	if !success {
		timeCache = time.Minute * 5
	}

	err = r.client.Set(ctx, key, bData, timeCache).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisCache) Del(ctx context.Context, key string, options ...cache.OptionCache) error {

	o := r.options.Clone()
	o.Load(options)

	if !checkIsKeyPattern(o) {
		err := r.client.Del(ctx, key).Err()
		if err != nil {
			return err
		}
	}
	keys, err := r.client.Keys(ctx, key).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	for _, key := range keys {
		err := r._del(ctx, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r redisCache) _del(ctx context.Context, key string) error {
	logger.Tracef("Redis cache delete %s %s", key)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
