package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"notifications/pkg/util/strset"
)

type SlidingWindowRateLimiter struct {
	slidingWindow *limiter
}

func (l *limiter) NewSlidingWindowLimiter(key string, rate float64, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		slidingWindow: &limiter{
			cache:     l.cache,
			keyPrefix: ":" + key,
			rate:      rate,
			window:    window,
		},
	}
}

func (sw *SlidingWindowRateLimiter) IsAllowed() bool {
	err := sw.slidingWindow.increment()
	if err != nil {
		return false
	}

	err = sw.slidingWindow.removeExpired()
	if err != nil {
		return false
	}

	count, err := sw.slidingWindow.countRequest()
	if err != nil {
		return false
	}

	allowedRequests := int64(sw.slidingWindow.rate * sw.slidingWindow.window.Seconds())

	return count <= allowedRequests
}

// 1e9 is a floating-point number that represents a billion using scientific notation, 1 * 10^9
// "+inf" is a string that represents positive infinity
const (
	_billion  = 1e9
	_infinity = "+inf"
)

func (l *limiter) increment() error {
	var (
		now    = time.Now().UnixNano()
		rate   = float64(now)
		member = strset.IntToStr(int(now))
	)

	_, err := l.cache.ZAdd(context.Background(), l.keyPrefix, rate, member)
	if err != nil {
		return err
	}
	return nil
}

func (l *limiter) removeExpired() error {
	var (
		now      = time.Now().UnixNano()
		minScore = float64(now) - l.window.Seconds()*_billion
	)

	_, err := l.cache.ZRemRangeByScore(context.Background(), l.keyPrefix, "0", fmt.Sprintf("%.0f", minScore))
	if err != nil {
		return err
	}
	return nil
}

func (l *limiter) countRequest() (int64, error) {
	var (
		now      = time.Now().UnixNano()
		minScore = float64(now) - l.window.Seconds()*_billion
	)

	count, err := l.cache.ZCount(context.Background(), l.keyPrefix, fmt.Sprintf("%.0f", minScore), _infinity)
	if err != nil {
		return 0, err
	}
	return count, nil
}
