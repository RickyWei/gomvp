package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
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
		err = errors.Wrap(err, "")
		return
	}
	err = json.Unmarshal(bs, &data)
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}

	return
}

func SetNx[T any](ctx context.Context, id string, exp time.Duration, data T) (err error) {
	key := getKey[T](id)
	bs, err := json.Marshal(data)
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}

	err = redisClient.SetNX(ctx, key, bs, exp).Err()
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}

	return
}
