package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_maxConn                  = 50
	_defaultConnAttempts      = 3
	_defaultConnTimeout       = time.Second
	_defaultHealthCheckPeriod = time.Second
	_defaultIdleTimeout       = 5 * time.Minute
)

type Postgres struct {
	maxConn      int
	connAttempts int
	connTimeout  time.Duration
	Pool         *pgxpool.Pool
}

func New(appName, url string, opts ...Option) (*Postgres, error) {
	var pg = &Postgres{
		maxConn:      _maxConn,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, errors.New("postgres: err on pgxpool.ParseConfig - " + err.Error())
	}

	cfg.MaxConns = int32(pg.maxConn)
	cfg.HealthCheckPeriod = _defaultHealthCheckPeriod
	cfg.MaxConnIdleTime = _defaultIdleTimeout
	cfg.ConnConfig.RuntimeParams["application_name"] = appName

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
		if err == nil {
			break
		}

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}
	if err != nil {
		return nil, errors.New("postgres: connAttempts = 0 - " + err.Error())
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

type Option func(*Postgres)

func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		if size != 0 {
			c.maxConn = size
		}
	}
}
