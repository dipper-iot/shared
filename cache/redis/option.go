package redis

import "gitlab.com/dipper-iot/shared/cache"

const keyIsKeyPattern = "isKeyPatternRedis"

func IsKeyPattern() cache.OptionCache {
	return func(optionsCache *cache.OptionsCache) {
		optionsCache.SetMeta(keyIsKeyPattern, true)
	}
}

func checkIsKeyPattern(o *cache.OptionsCache) bool {
	is, success := o.GetMeta(keyIsKeyPattern)
	if !success {
		return false
	}
	return is.(bool)
}
