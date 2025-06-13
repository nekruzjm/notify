package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"notifications/internal/repo/apiclient"
	"notifications/internal/service/admin"
	"notifications/pkg/lib/cache"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/security/ratelimiter"
)

var Module = fx.Provide(New)

type Protector interface {
	ProtectExternal() gin.HandlerFunc
	ProtectInternal() gin.HandlerFunc
}

type Params struct {
	fx.In

	Logger        logger.Logger
	RateLimiter   ratelimiter.Limiter
	Cache         cache.Cache
	APIClientRepo apiclient.Repo
	Admin         admin.Service
}

type mw struct {
	logger        logger.Logger
	rateLimiter   ratelimiter.Limiter
	cache         cache.Cache
	apiClientRepo apiclient.Repo
	admin         admin.Service
}

func New(p Params) Protector {
	return &mw{
		logger:        p.Logger,
		rateLimiter:   p.RateLimiter,
		cache:         p.Cache,
		apiClientRepo: p.APIClientRepo,
		admin:         p.Admin,
	}
}
