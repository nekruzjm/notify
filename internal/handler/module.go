package handler

import (
	"go.uber.org/fx"

	"notifications/internal/handler/broker"
	"notifications/internal/handler/http"
)

var Module = fx.Options(
	http.Module,
	broker.Module,
)
