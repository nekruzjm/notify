package broker

import (
	"notifications/internal/api/transport/broker/consumer"
	"notifications/internal/api/transport/broker/stream"
	"notifications/internal/api/transport/broker/subject"
)

func (p Params) registerJobs() {
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsJobEventRun, consumer.NotificationsJobEventRunProcessor, p.Event.Run)
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsJobPushCleaned, consumer.NotificationsJobPushCleanProcessor, p.Push.Clean)
}
