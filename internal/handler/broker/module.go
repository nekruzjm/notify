package broker

import (
	"go.uber.org/fx"

	"notifications/internal/handler/broker/email"
	"notifications/internal/handler/broker/event"
	"notifications/internal/handler/broker/push"
	"notifications/internal/handler/broker/sms"
	"notifications/internal/handler/broker/telegram"
	"notifications/internal/handler/broker/user"
)

var Module = fx.Options(
	push.Module,
	event.Module,
	user.Module,
	email.Module,
	telegram.Module,
	sms.Module,
)
