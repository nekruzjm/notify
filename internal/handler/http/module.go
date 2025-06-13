package http

import (
	"go.uber.org/fx"

	"notifications/internal/handler/http/event"
	"notifications/internal/handler/http/push"
)

var Module = fx.Options(
	push.Module,
	event.Module,
)
