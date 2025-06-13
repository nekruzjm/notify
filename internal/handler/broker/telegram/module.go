package telegram

import (
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"

	"notifications/internal/service/telegram"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Handler interface {
	Sent(jetstream.Msg)
}

type Params struct {
	fx.In

	Logger  logger.Logger
	Service telegram.Service
}

type handler struct {
	logger  logger.Logger
	service telegram.Service
}

func New(p Params) Handler {
	return &handler{
		logger:  p.Logger,
		service: p.Service,
	}
}
