package push

import (
	"context"
	"crypto/sha1"
	"encoding/binary"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"notifications/internal/repo/push"
	"notifications/internal/repo/rom"
	"notifications/internal/repo/user"
	"notifications/pkg/lib/broker/nats"
	"notifications/pkg/lib/cache"
	"notifications/pkg/lib/config"
	"notifications/pkg/lib/notifier/firebase"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
)

var Module = fx.Provide(New)

type Service interface {
	sender
	cleaner
}

type sender interface {
	Send(context.Context, *Request) (string, error)
}

type cleaner interface {
	Clean()
}

type Params struct {
	fx.In

	Config    config.Config
	Logger    logger.Logger
	Sentry    sentry.Sentry
	Nats      nats.Event
	Cache     cache.Cache
	FcmSender firebase.Sender
	UserRepo  user.Repo
	PushRepo  push.Repo
	RomRepo   rom.Repo
}

type service struct {
	channel map[bool]Service
}

func New(p Params) Service {
	hash := sha1.Sum([]byte(p.Config.GetString("podName")))
	podID := int64(binary.LittleEndian.Uint64(hash[:8]) % 1024)
	if podID == 0 {
		podID = 1
	}

	idGenerator, err := snowflake.NewNode(podID)
	if err != nil {
		p.Logger.Error("snowflake.NewNode", zap.Error(err))
		return nil
	}

	return &service{
		channel: map[bool]Service{
			true: &internal{
				logger:      p.Logger,
				sentry:      p.Sentry,
				nats:        p.Nats,
				cache:       p.Cache,
				fcmSender:   p.FcmSender,
				userRepo:    p.UserRepo,
				pushRepo:    p.PushRepo,
				romRepo:     p.RomRepo,
				idGenerator: idGenerator,
			},
			false: &external{
				logger:      p.Logger,
				sentry:      p.Sentry,
				nats:        p.Nats,
				fcmSender:   p.FcmSender,
				userRepo:    p.UserRepo,
				pushRepo:    p.PushRepo,
				romRepo:     p.RomRepo,
				idGenerator: idGenerator,
			},
		},
	}
}

func (s *service) Send(ctx context.Context, request *Request) (string, error) {
	return s.channel[request.IsInternal].Send(ctx, request)
}

func (s *service) Clean() {
	s.channel[true].Clean()
}
