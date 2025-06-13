package job

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"notifications/internal/api/transport/broker/stream"
	"notifications/internal/api/transport/broker/subject"
	"notifications/pkg/lib/broker/nats"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/scheduler"
)

var Module = fx.Invoke(New)

type Params struct {
	fx.In
	fx.Lifecycle

	Scheduler scheduler.Scheduler
	Nats      nats.Event
	Logger    logger.Logger
}

func New(p Params) {
	_, _ = p.Scheduler.Every(1).Minute().Do(p.launchEventRunner)
	_, _ = p.Scheduler.Every(60 * 24).Minute().Do(p.launchPushCleaner)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			p.Logger.Info("Notification worker started")
			return nil
		},
		OnStop: func(_ context.Context) error {
			p.Logger.Info("Notification worker stopped")
			return nil
		},
	})
}

func (p Params) launchEventRunner() {
	if err := p.Nats.Publish(stream.Notifications, subject.NotificationsJobEventRun, nil); err != nil {
		p.Logger.Error("err publishing event run", zap.Error(err))
	}
}

func (p Params) launchPushCleaner() {
	if err := p.Nats.Publish(stream.Notifications, subject.NotificationsJobPushCleaned, nil); err != nil {
		p.Logger.Error("err publishing push cleaned", zap.Error(err))
	}
}
