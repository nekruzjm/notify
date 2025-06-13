package user

import (
	"context"

	"go.uber.org/fx"

	"notifications/internal/repo/event"
	"notifications/internal/repo/user"
	"notifications/pkg/lib/notifier/firebase"
	"notifications/pkg/lib/observer/logger"
	"notifications/pkg/lib/observer/sentry"
)

var Module = fx.Provide(New)

type Service interface {
	CreateUser(context.Context, User) error
	UpdateToken(ctx context.Context, userID int, token string) error
	UpdateUserSettings(ctx context.Context, userID int, language string, isEnabled *bool) error
	UpdateStatus(ctx context.Context, userID int, status string) error
	UpdatePhone(ctx context.Context, userID int, phone string) error
	UpdatePersonExternalRef(ctx context.Context, userID int, personExternalRef string) error
}

type Params struct {
	fx.In

	Logger      logger.Logger
	Sentry      sentry.Sentry
	FcmTopicMan firebase.TopicManager
	UserRepo    user.Repo
	EventRepo   event.Repo
}

type service struct {
	logger      logger.Logger
	sentry      sentry.Sentry
	fcmTopicMan firebase.TopicManager
	userRepo    user.Repo
	eventRepo   event.Repo
}

func New(p Params) Service {
	return &service{
		logger:      p.Logger,
		sentry:      p.Sentry,
		userRepo:    p.UserRepo,
		eventRepo:   p.EventRepo,
		fcmTopicMan: p.FcmTopicMan,
	}
}
