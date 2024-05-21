package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/logger"
)

func getKey[T any](id string) string {
	var zero [0]T
	t := reflect.TypeOf(zero).Elem()

	return fmt.Sprintf(fmtKey, t.Name(), id)
}

func Get[T any](ctx context.Context, id string) (data T, err error) {
	key := getKey[T](id)
	bs, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		logger.L().Debug("redis get failed", zap.String("key", key), zap.Error(err))
		return
	}
	err = json.Unmarshal(bs, &data)
	if err != nil {
		logger.L().Debug("redis unmarshal failed", zap.String("bs", string(bs)), zap.Error(err))
		return
	}

	return
}

func SetNx[T any](ctx context.Context, id string, exp time.Duration, data T) (err error) {
	key := getKey[T](id)
	bs, err := json.Marshal(data)
	if err != nil {
		logger.L().Debug("redis marshal failed", zap.Any("data", data), zap.Error(err))
		return
	}

	err = redisClient.SetNX(ctx, key, bs, exp).Err()
	if err != nil {
		logger.L().Debug("redis setnx failed", zap.String("key", key), zap.Error(err))
		return
	}

	return
}
