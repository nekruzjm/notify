package push

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"

	"notifications/internal/service/push"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Handler interface {
	Sent(jetstream.Msg)
	Clean(jetstream.Msg)
	SyncSent(*nats.Msg)
}

type Params struct {
	fx.In

	Logger  logger.Logger
	Service push.Service
}

type handler struct {
	logger  logger.Logger
	service push.Service
}

func New(p Params) Handler {
	return &handler{
		logger:  p.Logger,
		service: p.Service,
	}
}
