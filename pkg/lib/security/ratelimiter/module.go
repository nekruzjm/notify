package ratelimiter

import (
	"time"

	"go.uber.org/fx"

	"notifications/pkg/lib/cache"
)

var Module = fx.Provide(New)

type Limiter interface {
	NewSlidingWindowLimiter(key string, rate float64, window time.Duration) *SlidingWindowRateLimiter
}

type Params struct {
	fx.In

	Cache cache.Cache
}

type limiter struct {
	cache cache.Cache

	keyPrefix string
	rate      float64
	window    time.Duration
}

func New(p Params) Limiter {
	return &limiter{
		cache: p.Cache,
	}
}
