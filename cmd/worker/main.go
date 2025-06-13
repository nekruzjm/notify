package main

import (
	"go.uber.org/fx"

	"notifications/cmd/worker/job"
	"notifications/pkg/lib/broker/nats"
	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/scheduler"
)

func main() {
	fx.New(
		job.Module,
		nats.Module,
		logger.Module,
		config.WorkerModule,
		scheduler.Module,
	).Run()
}
