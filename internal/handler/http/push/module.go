package push

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"notifications/internal/service/push"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Handler interface {
	Send(*gin.Context)
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
