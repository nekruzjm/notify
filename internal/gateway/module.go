package gateway

import (
	"go.uber.org/fx"

	"notifications/internal/gateway/admin"
)

var Module = fx.Options(
	admin.Module,
)
