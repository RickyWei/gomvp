package middleware

import (
	"context"
	"fmt"

	redisCache "github.com/go-redis/cache/v9"
	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/cache"
	"github.com/rickywei/sparrow/project/logger"
)

var (
	ac = &ApqCache{}
)

func GetApqCache() *ApqCache {
	return ac
}

type ApqCache struct{}

func (c *ApqCache) Add(ctx context.Context, key string, value any) {
	if err := cache.Cache.Set(&redisCache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf(cache.ApqKeyFmt, key),
		Value: value,
	}); err != nil {
		logger.L().Error("apq add failed", zap.String("key", key), zap.Error(err))
	}
}

func (c *ApqCache) Get(ctx context.Context, key string) (value any, ok bool) {
	if err := cache.Cache.Get(ctx, fmt.Sprintf(cache.ApqKeyFmt, key), &value); err != nil {
		logger.L().Error("apq get failed", zap.String("key", key), zap.Error(err))
		return
	}

	return
}
