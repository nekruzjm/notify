package user

import (
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"

	"notifications/internal/service/user"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Handler interface {
	UserCreated(jetstream.Msg)
	TokenUpdated(jetstream.Msg)
	StatusUpdated(jetstream.Msg)
	SettingsUpdated(jetstream.Msg)
	PhoneUpdated(jetstream.Msg)
	PersonExternalRefUpdated(jetstream.Msg)
}

type Params struct {
	fx.In

	Logger  logger.Logger
	Service user.Service
}

type handler struct {
	logger  logger.Logger
	service user.Service
}

func New(p Params) Handler {
	return &handler{
		logger:  p.Logger,
		service: p.Service,
	}
}
