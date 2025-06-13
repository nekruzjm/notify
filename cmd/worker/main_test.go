package main

import (
	"testing"

	"go.uber.org/fx"

	"notifications/cmd/worker/job"
	"notifications/pkg/lib/broker/nats"
	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/scheduler"
)

func Test_Deps(t *testing.T) {
	if err := fx.ValidateApp(deps()); err != nil {
		t.Error("err occurred during dependency injection:", err)
		return
	}
}

func deps() fx.Option {
	return fx.Options(
		job.Module,
		nats.Module,
		logger.Module,
		config.WorkerModule,
		scheduler.Module,
	)
}
