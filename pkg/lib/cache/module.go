package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Params struct {
	fx.In

	Config config.Config
	Logger logger.Logger
}

type cache struct {
	config config.Config
	logger logger.Logger

	client        *redis.Client
	clientCluster *redis.ClusterClient
	isCluster     bool
}

type Cache interface {
	writer
	reader
	pipeliner
}

type writer interface {
	Set(ctx context.Context, key string, value any, dur time.Duration) error
	SetObj(ctx context.Context, key string, value any, dur time.Duration) error
	ZAdd(ctx context.Context, key string, score float64, member any) (int64, error)
	ZRemRangeByScore(ctx context.Context, key string, minScore, maxScore string) (int64, error)
	Delete(ctx context.Context, key string) error
}

type reader interface {
	Get(ctx context.Context, key string, value any) error
	ZCount(ctx context.Context, key string, minScore, maxScore string) (int64, error)
	Exists(ctx context.Context, key string) bool
}

type pipeliner interface {
	Pipeline() redis.Pipeliner
}

const (
	_defaultPoolSize     = 100
	_defaultPoolTimeout  = 2 * time.Minute
	_defaultDialTimeout  = 5 * time.Minute
	_defaultReadTimeout  = 2 * time.Minute
	_defaultWriteTimeout = 2 * time.Minute
)

func New(p Params) Cache {
	var (
		cluster = p.Config.GetBool("redis.cluster")
		url     = p.Config.GetString("redis.url")
	)

	if cluster {
		var (
			ports     = p.Config.GetStringSlice("redis.ports")
			addresses = make([]string, 0, len(ports))
		)

		for _, port := range ports {
			addresses = append(addresses, url+":"+port)
		}

		clientCluster := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        addresses,
			PoolSize:     _defaultPoolSize,
			PoolTimeout:  _defaultPoolTimeout,
			DialTimeout:  _defaultDialTimeout,
			ReadTimeout:  _defaultReadTimeout,
			WriteTimeout: _defaultWriteTimeout,
		})

		clientCluster.Ping(context.Background())

		return &cache{
			isCluster:     true,
			clientCluster: clientCluster,
			logger:        p.Logger,
			config:        p.Config,
		}
	}

	var (
		port     = p.Config.GetString("redis.port")
		password = p.Config.GetString("redis.password")
		db       = p.Config.GetInt("redis.db")
	)

	client := redis.NewClient(&redis.Options{
		Addr:     url + ":" + port,
		Password: password,
		DB:       db,
	})

	client.Ping(context.Background())

	return &cache{
		isCluster: false,
		client:    client,
		logger:    p.Logger,
		config:    p.Config,
	}
}
