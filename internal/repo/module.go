package repo

import (
	"go.uber.org/fx"

	"notifications/internal/repo/apiclient"
	"notifications/internal/repo/event"
	"notifications/internal/repo/push"
	"notifications/internal/repo/rom"
	"notifications/internal/repo/user"
)

var Module = fx.Options(
	user.Module,
	push.Module,
	apiclient.Module,
	event.Module,
	rom.Module,
)
