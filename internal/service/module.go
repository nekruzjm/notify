package service

import (
	"go.uber.org/fx"

	"notifications/internal/service/admin"
	"notifications/internal/service/email"
	"notifications/internal/service/event"
	"notifications/internal/service/push"
	"notifications/internal/service/sms"
	"notifications/internal/service/telegram"
	"notifications/internal/service/user"
)

var Module = fx.Options(
	admin.Module,
	push.Module,
	user.Module,
	event.Module,
	email.Module,
	telegram.Module,
	sms.Module,
)
