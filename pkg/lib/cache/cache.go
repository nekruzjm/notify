package cache

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

const (
	_defaultServicePrefix = "notifications"
)

func (c *cache) Set(ctx context.Context, key string, value any, dur time.Duration) error {
	if c.isCluster {
		return c.clientCluster.Set(ctx, _defaultServicePrefix+key, value, dur).Err()
	}
	return c.client.Set(ctx, _defaultServicePrefix+key, value, dur).Err()
}

func (c *cache) SetObj(ctx context.Context, key string, value any, dur time.Duration) error {
	bytes, err := sonic.Marshal(value)
	if err != nil {
		return err
	}

	if c.isCluster {
		return c.clientCluster.Set(ctx, _defaultServicePrefix+key, bytes, dur).Err()
	}
	return c.client.Set(ctx, _defaultServicePrefix+key, bytes, dur).Err()
}

func (c *cache) Exists(ctx context.Context, key string) bool {
	if c.isCluster {
		exist, _ := c.clientCluster.Exists(ctx, _defaultServicePrefix+key).Result()
		return exist > 0
	}
	exist, _ := c.client.Exists(ctx, _defaultServicePrefix+key).Result()
	return exist > 0
}

func (c *cache) ZAdd(ctx context.Context, key string, score float64, member any) (int64, error) {
	if c.isCluster {
		return c.clientCluster.ZAdd(ctx, _defaultServicePrefix+key, redis.Z{
			Score:  score,
			Member: member,
		}).Result()
	}
	return c.client.ZAdd(ctx, _defaultServicePrefix+key, redis.Z{
		Score:  score,
		Member: member,
	}).Result()
}

func (c *cache) ZRemRangeByScore(ctx context.Context, key string, minScore, maxScore string) (int64, error) {
	if c.isCluster {
		return c.clientCluster.ZRemRangeByScore(ctx, _defaultServicePrefix+key, minScore, maxScore).Result()
	}
	return c.client.ZRemRangeByScore(ctx, _defaultServicePrefix+key, minScore, maxScore).Result()
}

func (c *cache) ZCount(ctx context.Context, key string, minScore, maxScore string) (int64, error) {
	if c.isCluster {
		return c.clientCluster.ZCount(ctx, _defaultServicePrefix+key, minScore, maxScore).Result()
	}
	return c.client.ZCount(ctx, _defaultServicePrefix+key, minScore, maxScore).Result()
}

func (c *cache) Get(ctx context.Context, key string, value any) error {
	if c.isCluster {
		val, err := c.clientCluster.Get(ctx, _defaultServicePrefix+key).Result()
		if err != nil {
			return err
		}
		return sonic.Unmarshal([]byte(val), value)
	}

	val, err := c.client.Get(ctx, _defaultServicePrefix+key).Result()
	if err != nil {
		return err
	}
	return sonic.Unmarshal([]byte(val), value)
}

func (c *cache) Delete(ctx context.Context, key string) error {
	if c.isCluster {
		return c.clientCluster.Del(ctx, _defaultServicePrefix+key).Err()
	}
	return c.client.Del(ctx, _defaultServicePrefix+key).Err()
}

func (c *cache) Pipeline() redis.Pipeliner {
	if c.isCluster {
		return c.clientCluster.Pipeline()
	}
	return c.client.Pipeline()
}
