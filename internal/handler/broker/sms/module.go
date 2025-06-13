package sms

import (
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"

	"notifications/internal/service/sms"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Handler interface {
	Sent(jetstream.Msg)
}

type Params struct {
	fx.In

	Logger  logger.Logger
	Service sms.Service
}

type handler struct {
	logger  logger.Logger
	service sms.Service
}

func New(p Params) Handler {
	return &handler{
		logger:  p.Logger,
		service: p.Service,
	}
}
