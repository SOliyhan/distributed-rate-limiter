package limiter

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed scripts/token_bucket.lua
var tokenBucketLua string

type RedisLimiter struct {
	client   *redis.Client
	capacity int
	rate     float64
}

func NewRedisLimiter(client *redis.Client, capacity int, rate float64) *RedisLimiter {
	return &RedisLimiter{
		client:   client,
		capacity: capacity,
		rate:     rate,
	}
}

func (rl *RedisLimiter) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now().Unix()

	res, err := rl.client.Eval(
		ctx,
		tokenBucketLua,
		[]string{fmt.Sprintf("rl:%s", key)},
		rl.capacity,
		rl.rate,
		now,
	).Int()

	if err != nil {
		return false, err
	}

	return res == 1, nil
}
