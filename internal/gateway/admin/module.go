package admin

import (
	"context"

	"go.uber.org/fx"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Gateway interface {
	Authorize(ctx context.Context, token string) (Admin, error)
}

type Params struct {
	fx.In

	Config config.Config
	Logger logger.Logger
}

type gateway struct {
	config config.Config
	logger logger.Logger
}

func New(p Params) Gateway {
	return &gateway{
		config: p.Config,
		logger: p.Logger,
	}
}
