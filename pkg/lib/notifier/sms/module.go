package sms

import (
	"context"

	"go.uber.org/fx"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type SMS interface {
	Send(ctx context.Context, request Request) error
}

type Params struct {
	fx.In

	Config config.Config
	Logger logger.Logger
}

type sms struct {
	config config.Config
	logger logger.Logger
}

func New(p Params) SMS {
	return &sms{
		config: p.Config,
		logger: p.Logger,
	}
}
