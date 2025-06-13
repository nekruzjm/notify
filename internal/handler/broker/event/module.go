package event

import (
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"

	"notifications/internal/service/event"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Handler interface {
	Run(jetstream.Msg)
	TopicSubscribed(jetstream.Msg)
	TopicUnsubscribed(jetstream.Msg)
}

type Params struct {
	fx.In

	Logger  logger.Logger
	Service event.Service
}

type handler struct {
	logger  logger.Logger
	service event.Service
}

func New(p Params) Handler {
	return &handler{
		logger:  p.Logger,
		service: p.Service,
	}
}
