package admin

import (
	"context"

	"go.uber.org/fx"

	"notifications/internal/gateway/admin"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
)

var Module = fx.Provide(New)

type Service interface {
	Authorize(ctx context.Context, token string) (Admin, error)
}

type Params struct {
	fx.In

	Logger  logger.Logger
	Sentry  sentry.Sentry
	Gateway admin.Gateway
}

type service struct {
	logger  logger.Logger
	sentry  sentry.Sentry
	gateway admin.Gateway
}

func New(p Params) Service {
	return &service{
		logger:  p.Logger,
		sentry:  p.Sentry,
		gateway: p.Gateway,
	}
}
