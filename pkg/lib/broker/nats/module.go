package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Event interface {
	Publish(stream, subj string, msg any) error
	Subscribe(stream, subj, consumer string, handler jetstream.MessageHandler, opts ...SubscriptionOptions)
	Reply(subj, qGroup string, handler nats.MsgHandler)
}

type Params struct {
	fx.In
	fx.Lifecycle

	Logger logger.Logger
	Config config.Config
}

type natsConn struct {
	logger logger.Logger
	config config.Config
	ctx    context.Context
	js     jetstream.JetStream
	conn   *nats.Conn

	replicas int
}

func New(p Params) Event {
	var nc = &natsConn{
		logger: p.Logger,
		config: p.Config,
		ctx:    nats.Context(context.Background()),
	}

	conn, err := nc._connect()
	if err != nil {
		return nil
	}

	p.Lifecycle.Append(fx.StopHook(func() { conn.Close() }))

	return nc
}
