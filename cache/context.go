package cache

import (
	"context"
	"time"
)

type keyTime struct{}
type keyTimeout struct{}

var keyTimeData = keyTime{}
var keyTimoutData = keyTimeout{}

func TimeCacheToContext(ctx context.Context, timeCache time.Duration) context.Context {
	return context.WithValue(ctx, keyTimeData, timeCache)
}

func TimeCacheFromContext(ctx context.Context) (time.Duration, bool) {
	data := ctx.Value(keyTimeData)
	if data != nil {
		return data.(time.Duration), true
	}

	return 0, false
}

func TimeoutToContext(ctx context.Context, timeCache time.Duration) context.Context {
	return context.WithValue(ctx, keyTimoutData, timeCache)
}

func TimeoutFromContext(ctx context.Context) (time.Duration, bool) {
	data := ctx.Value(keyTimoutData)
	if data != nil {
		return data.(time.Duration), true
	}

	return 0, false
}
