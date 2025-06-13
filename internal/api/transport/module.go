package transport

import (
	"go.uber.org/fx"

	"notifications/internal/api/transport/broker"
	"notifications/internal/api/transport/http"
	"notifications/internal/api/transport/http/middleware"
)

var Module = fx.Options(
	http.Module,
	broker.Module,
	middleware.Module,
)
