package cache

import "context"

type Cache interface {
	Init(options ...OptionCache) error
	Get(ctx context.Context, key string, data interface{}, options ...OptionCache) (bool, error)
	Set(ctx context.Context, key string, data interface{}, options ...OptionCache) error
	Del(ctx context.Context, key string, options ...OptionCache) error
}
