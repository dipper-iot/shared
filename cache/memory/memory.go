package memory

import (
	"context"
	"gitlab.com/dipper-iot/shared/cache"
	"gitlab.com/dipper-iot/shared/util"
	"sync"
)

type Cache struct {
	mapData *sync.Map
}

func NewCache() *Cache {
	return &Cache{
		mapData: &sync.Map{},
	}
}

func (c Cache) Init(options ...cache.OptionCache) error {

	return nil
}

func (c Cache) Get(ctx context.Context, key string, data interface{}, options ...cache.OptionCache) (bool, error) {
	raw, ok := c.mapData.Load(key)
	if ok {
		return true, util.Mapper(raw, data)
	}
	return false, nil
}

func (c Cache) Set(ctx context.Context, key string, data interface{}, options ...cache.OptionCache) error {
	c.mapData.Store(key, data)
	return nil
}

func (c Cache) Del(ctx context.Context, key string, options ...cache.OptionCache) error {
	c.mapData.Delete(key)
	return nil
}
