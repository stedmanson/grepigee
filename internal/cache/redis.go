package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func InitRedis(addr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(ctx, key, value, expiration).Err()
}
