package sentry

import (
	"context"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Sentry interface {
	CaptureException(error)
	CurrentHub() *sentry.Hub
}

type Params struct {
	fx.In
	fx.Lifecycle

	Config config.Config
	Logger logger.Logger
}

type sentryObserver struct {
	hub    *sentry.Hub
	logger logger.Logger
	stage  string
}

func New(p Params) Sentry {
	var options = sentry.ClientOptions{
		Dsn:              p.Config.GetString("sentry.dsn"),
		Debug:            p.Config.GetBool("sentry.debug"),
		ServerName:       p.Config.GetString("sentry.serverName"),
		Environment:      p.Config.GetString("sentry.environment"),
		EnableTracing:    true,
		AttachStacktrace: true,
		TracesSampleRate: 1.,
	}

	p.Logger.Info(p.Config.GetString("sentry.dsn"))

	client, err := sentry.NewClient(options)
	if err != nil {
		p.Logger.Error("can't run Sentry", zap.Error(err))
		os.Exit(1)
	}

	hub := sentry.CurrentHub()
	hub.BindClient(client)

	p.Lifecycle.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				p.Logger.Info("Sentry started")
				return nil
			},
			OnStop: func(_ context.Context) error {
				_ = hub.Flush(2 * time.Second)
				p.Logger.Info("Sentry stopped")
				return nil
			},
		},
	)

	return &sentryObserver{
		stage:  p.Config.GetString("sentry.stage"),
		hub:    hub,
		logger: p.Logger,
	}
}
