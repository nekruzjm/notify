package email

import (
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"

	"notifications/internal/service/email"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Handler interface {
	Sent(jetstream.Msg)
}

type Params struct {
	fx.In

	Logger  logger.Logger
	Service email.Service
}

type handler struct {
	logger  logger.Logger
	service email.Service
}

func New(p Params) Handler {
	return &handler{
		logger:  p.Logger,
		service: p.Service,
	}
}
