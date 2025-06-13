package scheduler

import "github.com/go-co-op/gocron"

func (s *scheduler) Every(interval any) *gocron.Scheduler {
	return s.scheduler.Every(interval)
}

func (s *scheduler) Cron(cronExpression string) *gocron.Scheduler {
	return s.scheduler.Cron(cronExpression)
}

func (s *scheduler) CronWithSeconds(cronExpression string) *gocron.Scheduler {
	return s.scheduler.CronWithSeconds(cronExpression)
}

func (s *scheduler) RemoveByTag(tag string) error {
	return s.scheduler.RemoveByTag(tag)
}
