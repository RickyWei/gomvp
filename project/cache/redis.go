package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/conf"
	"github.com/rickywei/sparrow/project/logger"
)



var (
	redisClient    redis.UniversalClient
	Cache *cache.Cache
)

func init() {
	var err error
	if viper.IsSet("redis.cluster") {
		redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: conf.Strings("redis.cluster.addrs"),
		})
	} else if viper.IsSet("redis.sentinel") {
		redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:      conf.Strings("redis.sentinel.addrs"),
			MasterName: conf.String("redis.sentinel.addrs"),
		})
	} else {
		redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: conf.Strings("redis.client.addr"),
		})
	}
	if _, err = redisClient.Ping(context.Background()).Result(); err != nil {
		logger.L().Fatal("ping redis failed", zap.Error(err))
	}

	Cache = cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: cache.NewTinyLFU(10000, time.Minute),
	})
}
