package rs

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/tinrab/retry"
	"os"
	"time"
)

func NewRedisClient(url string, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
	})

	_, err := client.Ping(context.TODO()).Result()
	if !errors.Is(err, nil) {
		return nil, err
	}

	retry.ForeverSleep(10*time.Second, func(_ int) error {
		_, err := client.Ping(context.TODO()).Result()
		return err
	})

	return client, nil
}

func NewRedisToEnv() (*redis.Client, error) {
	redisUrl := os.Getenv("REDIS_ADDRESS")
	pass := os.Getenv("REDIS_PASS")

	return NewRedisClient(redisUrl, pass)
}
