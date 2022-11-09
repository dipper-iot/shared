package cache

import (
	"context"
	"gitlab.com/dipper-iot/shared/logger"
)

func UtilSetCache(c Cache, ctx context.Context, key string, data interface{}) {
	err := c.Set(ctx, key, data)
	if err != nil {
		logger.Error(err)
	}
}
