package telegram

import (
	"context"

	"go.uber.org/fx"

	"notifications/pkg/lib/notifier/telegram"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
)

var Module = fx.Provide(New)

type Service interface {
	Send(context.Context, Message) error
}

type Params struct {
	fx.In

	Logger logger.Logger
	Sentry sentry.Sentry
	Tg     telegram.Telegram
}

type service struct {
	logger logger.Logger
	sentry sentry.Sentry
	tg     telegram.Telegram
}

func New(p Params) Service {
	return &service{
		logger: p.Logger,
		sentry: p.Sentry,
		tg:     p.Tg,
	}
}
