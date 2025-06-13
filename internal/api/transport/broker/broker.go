package broker

import (
	"go.uber.org/fx"

	"notifications/internal/api/transport/broker/consumer"
	"notifications/internal/api/transport/broker/stream"
	"notifications/internal/api/transport/broker/subject"
	"notifications/internal/handler/broker/email"
	"notifications/internal/handler/broker/event"
	"notifications/internal/handler/broker/push"
	"notifications/internal/handler/broker/sms"
	"notifications/internal/handler/broker/telegram"
	"notifications/internal/handler/broker/user"
	"notifications/pkg/lib/broker/nats"
)

var Module = fx.Options(fx.Invoke(RegisterEvents))

type Params struct {
	fx.In

	Nats nats.Event

	User  user.Handler
	Push  push.Handler
	Event event.Handler
	Email email.Handler
	Sms   sms.Handler
	Tg    telegram.Handler
}

func RegisterEvents(p Params) {
	// user data manager
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsUserCreated, consumer.NotificationsUserProcessor, p.User.UserCreated)
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsFcmRegistrationTokenUpdated, consumer.NotificationsFcmTokenProcessor, p.User.TokenUpdated)
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsUserStatusUpdated, consumer.NotificationsUserStatusProcessor, p.User.StatusUpdated)
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsUserSettingsUpdated, consumer.NotificationsUserSettingsProcessor, p.User.SettingsUpdated)
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsUserPhoneUpdated, consumer.NotificationsUserPhoneProcessor, p.User.PhoneUpdated)
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsUserPersonRefUpdated, consumer.NotificationsUserPersonRefProcessor, p.User.PersonExternalRefUpdated)
	// notifier
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsPushSent, consumer.NotificationsPushProcessor, p.Push.Sent, nats.WithMaxDelivery(1))
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsEmailSent, consumer.NotificationsEmailProcessor, p.Email.Sent)
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsSmsSent, consumer.NotificationsSmsProcessor, p.Sms.Sent)
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsTgSent, consumer.NotificationsTgProcessor, p.Tg.Sent)
	// fcm topic subscribe/unsubscribe
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsTopicUsersSubscribed, consumer.NotificationsTopicUsersSubProcessor, p.Event.TopicSubscribed)
	p.Nats.Subscribe(stream.Notifications, subject.NotificationsTopicUsersUnsubscribed, consumer.NotificationsTopicUsersUnsubProcessor, p.Event.TopicUnsubscribed)

	p.Nats.Reply(subject.NotificationsSyncPushSent, consumer.NotificationsGroup, p.Push.SyncSent)

	p.registerJobs()
}
