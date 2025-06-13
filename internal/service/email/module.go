package email

import (
	"context"
	"net/smtp"

	"go.uber.org/fx"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
)

var Module = fx.Provide(New)

type Service interface {
	Send(context.Context, Email) error
}

type Params struct {
	fx.In

	Config config.Config
	Logger logger.Logger
	Sentry sentry.Sentry
}

type service struct {
	config    config.Config
	logger    logger.Logger
	sentry    sentry.Sentry
	plainAuth smtp.Auth
}

func New(p Params) Service {
	return &service{
		config: p.Config,
		logger: p.Logger,
		sentry: p.Sentry,
		plainAuth: smtp.PlainAuth("",
			p.Config.GetString("email.from"),
			p.Config.GetString("email.password"),
			p.Config.GetString("email.host"),
		),
	}
}
