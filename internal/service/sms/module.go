package sms

import (
	"context"

	"go.uber.org/fx"

	"notifications/pkg/lib/notifier/sms"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
)

var Module = fx.Provide(New)

type Service interface {
	Send(context.Context, map[string]string) error
}

type Params struct {
	fx.In

	Logger logger.Logger
	Sentry sentry.Sentry
	Sms    sms.SMS
}

type service struct {
	logger logger.Logger
	sentry sentry.Sentry
	sms    sms.SMS
}

func New(p Params) Service {
	return &service{
		logger: p.Logger,
		sentry: p.Sentry,
		sms:    p.Sms,
	}
}
