package redis_sync

import (
	"context"
	"gitlab.com/dipper-iot/shared/cache"
	"gitlab.com/dipper-iot/shared/logger"
	"os"
	"testing"
	"time"
)

var cacheData cache.Cache

func init() {
	var err error
	cacheData, err = NewCacheEnv()
	if err != nil {
		logger.Error(err)
		os.Exit(0)
	}
}

func TestDelCache(t *testing.T) {
	err := cacheData.Set(context.TODO(), "test-a-b-c-d", "test")
	if err != nil {
		t.Errorf("Set Cache is incorrect %v", err)
	}

	err = cacheData.Set(cache.TimeCacheToContext(context.TODO(), 5*time.Minute), "test-a-c", "test")
	if err != nil {
		t.Errorf("Set Cache is incorrect %v", err)
	}

	err = cacheData.Set(cache.TimeCacheToContext(context.TODO(), 5*time.Minute), "test-bd", "test")
	if err != nil {
		t.Errorf("Set Cache is incorrect %v", err)
	}

	err = cacheData.Del(context.TODO(), "test-*b*", IsKeyPattern())
	if err != nil {
		t.Errorf("Set Cache is incorrect %v", err)
	}

	var rs string
	success, err := cacheData.Get(context.TODO(), "test-a-b-c-d", &rs)
	if err != nil {
		t.Errorf("Get Cache is incorrect %v", err)
	}
	if success {
		t.Errorf("Delete Cache is not delete")
	}

	success, err = cacheData.Get(context.TODO(), "test-a-c", &rs)
	if err != nil {
		t.Errorf("Delete Cache is incorrect %v", err)
	}
	if !success {
		t.Errorf("Delete Cache is delete all")
	}
}
