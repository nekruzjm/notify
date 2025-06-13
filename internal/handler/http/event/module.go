package event

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"notifications/internal/service/event"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Handler interface {
	writer
	reader
	runner
	imageManager
}

type reader interface {
	Get(*gin.Context)
}

type writer interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type runner interface {
	LoadUsers(*gin.Context)
	LoadAllUsers(*gin.Context)
	Run(*gin.Context)
}

type imageManager interface {
	UploadImage(*gin.Context)
	RemoveImage(*gin.Context)
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
