package lib

import (
	"go.uber.org/fx"

	"notifications/pkg/lib/broker/nats"
	"notifications/pkg/lib/cache"
	"notifications/pkg/lib/config"
	"notifications/pkg/lib/fileman"
	"notifications/pkg/lib/notifier/firebase"
	"notifications/pkg/lib/notifier/sms"
	"notifications/pkg/lib/notifier/telegram"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
	"notifications/pkg/lib/scheduler"
	"notifications/pkg/lib/security/ratelimiter"
	"notifications/pkg/lib/tinypng"
)

var Module = fx.Options(
	config.Module,
	logger.Module,
	sentry.Module,
	scheduler.Module,
	nats.Module,
	cache.Module,
	ratelimiter.Module,
	fileman.Module,
	firebase.Module,
	telegram.Module,
	sms.Module,
	tinypng.Module,
)
