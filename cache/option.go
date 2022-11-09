package cache

import (
	"sync"
)

type OptionsCache struct {
	meta map[string]interface{}
	lock *sync.Mutex
}

func NewOptionsCache() *OptionsCache {
	return &OptionsCache{
		meta: map[string]interface{}{},
		lock: &sync.Mutex{},
	}
}

type OptionCache func(*OptionsCache)

func (o *OptionsCache) Clone() *OptionsCache {
	return &OptionsCache{
		meta: o.meta,
		lock: &sync.Mutex{},
	}
}

func (o *OptionsCache) SetMeta(key string, data interface{}) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.meta[key] = data
}

func (o *OptionsCache) GetMeta(key string) (interface{}, bool) {
	o.lock.Lock()
	defer o.lock.Unlock()
	result, success := o.meta[key]
	return result, success
}

func SetData(key string, data interface{}) OptionCache {
	return func(o *OptionsCache) {
		o.SetMeta(key, data)
	}
}

func (o *OptionsCache) Load(options []OptionCache) {
	if options == nil || len(options) == 0 {
		return
	}
	for _, option := range options {
		option(o)
	}
}
