package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Sender interface {
	SendPush(context.Context, *messaging.Message) (string, error)
	SendPushDryRun(ctx context.Context, message *messaging.Message) (string, error)
	SendEach(context.Context, []*messaging.Message) (*messaging.BatchResponse, error)
	SendEachDryRun(ctx context.Context, messages []*messaging.Message) (*messaging.BatchResponse, error)
	// SendMulticast sends a message to multiple devices like Android, Ios, Web in one request.
	SendMulticast(context.Context, *messaging.MulticastMessage) (*messaging.BatchResponse, error)
	SendMulticastDryRun(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
}

type TopicManager interface {
	// SubscribeTokens to a topic with a list of registration tokens.
	// Call SendPush to send a message to topic and all registered devices will get message.
	SubscribeTokens(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)
	UnsubscribeTokens(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)

	Subscribe(ctx context.Context, token, topic string) (*messaging.TopicManagementResponse, error)
	Unsubscribe(ctx context.Context, token, topic string) (*messaging.TopicManagementResponse, error)
}

type Params struct {
	fx.In

	Logger logger.Logger
}

type fb struct {
	logger logger.Logger
	client *messaging.Client
}

const (
	_projectID = "my-app"
	_authDir   = "./firebase.json"
)

func New(p Params) (Sender, TopicManager) {
	var (
		ctx = context.Background()
		cfg = &firebase.Config{ProjectID: _projectID}
	)

	app, err := firebase.NewApp(ctx, cfg, option.WithCredentialsFile(_authDir))
	if err != nil {
		p.Logger.Error("cannot create app instance", zap.Error(err))
		return nil, nil
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		p.Logger.Error("cannot create client instance", zap.Error(err))
		return nil, nil
	}

	p.Logger.Info("Connected to Firebase", zap.String("projectID", cfg.ProjectID))

	var fcm = &fb{
		logger: p.Logger,
		client: client,
	}

	return fcm, fcm
}
