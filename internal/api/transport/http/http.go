package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "notifications/docs"
	"notifications/internal/api/resp"
	"notifications/internal/api/resp/code"
	"notifications/internal/api/transport/http/middleware"
	"notifications/internal/handler/http/event"
	"notifications/internal/handler/http/push"
	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Options(fx.Invoke(NewHTTPRouter))

type Params struct {
	fx.In
	fx.Lifecycle

	Config     config.Config
	Logger     logger.Logger
	Middleware middleware.Protector

	Push  push.Handler
	Event event.Handler
}

// NewHTTPRouter
//
//	@Title						Notifications API
//	@Version					1.0
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.email				jamshedzodnekruz@gmail.com
//	@Host						api-notifications.dev.my.cloud
//	@BasePath					/api
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
func NewHTTPRouter(p Params) {
	gin.SetMode(gin.ReleaseMode)

	var (
		router       = gin.New()
		internalBase = router.Group("/api/notifications-internal/v1")
		externalBase = router.Group("/api/notifications-external/v1")
	)

	if p.Config.GetString("notifications.stage") == "dev" {
		router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	}

	router.Use(gin.Recovery())
	registerPprof(router)

	internalBase.GET("/health", func(c *gin.Context) { resp.JSON(c.Writer, code.Success, resp.Success) })

	internalEvents := internalBase.Group("/events").Use(p.Middleware.ProtectInternal())
	internalEvents.GET("/", p.Event.Get)
	internalEvents.POST("/", p.Event.Create)
	internalEvents.PUT("/:id", p.Event.Update)
	internalEvents.DELETE("/:id", p.Event.Delete)
	internalEvents.POST("/:id/load-users", p.Event.LoadUsers)
	internalEvents.POST("/:id/load-all-users", p.Event.LoadAllUsers)
	internalEvents.POST("/:id/run", p.Event.Run)
	internalEvents.POST("/:id/image/:language", p.Event.UploadImage)
	internalEvents.DELETE("/:id/image/:language", p.Event.RemoveImage)

	externalBase.Group("/push").Use(p.Middleware.ProtectExternal()).POST("/", p.Push.Send)

	var server = http.Server{
		Addr:    p.Config.GetString("notifications.server.port"),
		Handler: router.Handler(),
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					p.Logger.Error("err on server.ListenAndServe()", zap.Error(err), zap.String("server address", server.Addr))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_ = server.Shutdown(ctx)
			return nil
		},
	})
}
