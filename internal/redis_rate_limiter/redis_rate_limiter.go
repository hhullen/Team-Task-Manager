package ratelimiter

import (
	"context"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -destination=redis_rate_limiter_mock.go -package=ratelimiter . ILimiter

type ILimiter interface {
	Allow(ctx context.Context, key string, limit redis_rate.Limit) (*redis_rate.Result, error)
}

type RedisRateLimiter struct {
	ctx     context.Context
	limiter ILimiter
}

func NewRedisRateLimiter(ctx context.Context, r *redis.Client) *RedisRateLimiter {
	limiter := redis_rate.NewLimiter(r)
	return buildRedisRateLimiter(ctx, limiter)
}

func buildRedisRateLimiter(ctx context.Context, l ILimiter) *RedisRateLimiter {
	return &RedisRateLimiter{
		ctx:     ctx,
		limiter: l,
	}
}

func (rl *RedisRateLimiter) Allow(key string, perSecond int) (bool, error) {
	res, err := rl.limiter.Allow(rl.ctx, key, redis_rate.PerSecond(perSecond))
	if err != nil {
		return false, err
	}

	return res.Allowed > 0, nil
}
