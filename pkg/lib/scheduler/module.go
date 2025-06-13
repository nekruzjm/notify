package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Scheduler interface {
	Every(interval any) *gocron.Scheduler
	Cron(cronExpression string) *gocron.Scheduler
	CronWithSeconds(cronExpression string) *gocron.Scheduler
	RemoveByTag(tag string) error
}

type scheduler struct {
	fx.In

	scheduler *gocron.Scheduler
}

func New() Scheduler {
	schedule := gocron.NewScheduler(time.UTC)
	schedule.TagsUnique()
	schedule.StartAsync()

	return &scheduler{scheduler: schedule}
}
