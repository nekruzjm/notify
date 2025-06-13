package nats

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

const (
	_msgMaxAge     = 12 * time.Hour
	_maxWaitingMsg = 1024
	_maxDelivery   = 3
	_maxAckWait    = time.Second * 30
)

func (n *natsConn) Subscribe(stream, subj, consumer string, handler jetstream.MessageHandler, options ...SubscriptionOptions) {
	var opt = new(subscriptionOpt)
	err := n.configureOptions(opt, options...)
	if err != nil {
		n.logger.Error("err on configuring options ", zap.Error(err))
		return
	}

	n.subscribeWithOpt(
		withSubj(subj),
		withController(handler),
		withStreamCfg(&jetstream.StreamConfig{
			Name:        stream,
			Compression: jetstream.S2Compression,
			MaxAge:      opt.msgTTL,
			Storage:     jetstream.FileStorage,
			Retention:   jetstream.WorkQueuePolicy,
			Replicas:    n.replicas,
		}),
		withConsumerCfg(&jetstream.ConsumerConfig{
			Durable:       consumer,
			MaxWaiting:    _maxWaitingMsg,
			AckWait:       opt.ackWait,
			MaxDeliver:    opt.maxDelivery,
			AckPolicy:     jetstream.AckExplicitPolicy,
			ReplayPolicy:  jetstream.ReplayOriginalPolicy,
			DeliverPolicy: jetstream.DeliverAllPolicy,
			FilterSubject: subj,
			Replicas:      n.replicas,
		}),
	)
}

func (n *natsConn) subscribeWithOpt(options ...SubscriptionOptions) {
	if len(options) == 0 {
		n.logger.Error("no options were provided")
		return
	}

	opt := new(subscriptionOpt)
	err := n.configureOptions(opt, options...)
	if err != nil {
		n.logger.Error("err on configuring options", zap.Error(err))
		return
	}

	err = n.upsertStream(n.ctx, opt.subj, opt.streamCfg)
	if err != nil {
		n.logger.Error("can't upsert stream",
			zap.Error(err),
			zap.String("streamName", opt.streamCfg.Name),
			zap.String("subject", opt.subj),
			zap.String("consumer", opt.consumerCfg.Durable))
		return
	}

	con, err := n.js.CreateOrUpdateConsumer(n.ctx, opt.streamCfg.Name, *opt.consumerCfg)
	if err != nil {
		n.logger.Error("can't upsert consumer",
			zap.Error(err),
			zap.String("streamName", opt.streamCfg.Name),
			zap.String("subject", opt.subj),
			zap.String("consumer", opt.consumerCfg.Durable))
		return
	}

	_, err = con.Consume(opt.msgHandler)
	if err != nil {
		n.logger.Error("err on consume", zap.Error(err))
		return
	}

	n.logger.Info("Listening event",
		zap.String("streamName", opt.streamCfg.Name),
		zap.String("subject", opt.subj),
		zap.String("consumer", opt.consumerCfg.Durable))
}

func (n *natsConn) upsertStream(ctx context.Context, subj string, streamCfg *jetstream.StreamConfig) error {
	stream, err := n.js.Stream(ctx, streamCfg.Name)
	if err != nil {
		if !errors.Is(err, jetstream.ErrStreamNotFound) {
			return err
		}

		stream, err = n.js.CreateStream(ctx, *streamCfg)
		if err != nil {
			return err
		}
	}

	streamInfo, err := stream.Info(ctx)
	if err != nil {
		return err
	}

	streamInfo.Config.Compression = streamCfg.Compression
	streamInfo.Config.MaxAge = streamCfg.MaxAge
	streamInfo.Config.Storage = streamCfg.Storage
	streamInfo.Config.Retention = streamCfg.Retention
	streamInfo.Config.Replicas = streamCfg.Replicas

	if !slices.Contains(streamInfo.Config.Subjects, subj) {
		streamInfo.Config.Subjects = append(streamInfo.Config.Subjects, subj)
	}

	_, err = n.js.UpdateStream(ctx, streamInfo.Config)
	if err != nil {
		return err
	}

	return nil
}

func (n *natsConn) configureOptions(opt *subscriptionOpt, opts ...SubscriptionOptions) error {
	for _, o := range opts {
		if o != nil {
			if err := o.configureSubscription(opt); err != nil {
				return err
			}
		}
	}

	return nil
}

type SubscriptionOptions interface {
	configureSubscription(*subscriptionOpt) error
}

type subscriptionCfgFn func(*subscriptionOpt) error

func (opts subscriptionCfgFn) configureSubscription(sub *subscriptionOpt) error {
	return opts(sub)
}

type subscriptionOpt struct {
	subj        string
	maxDelivery int
	ackWait     time.Duration
	msgTTL      time.Duration
	msgHandler  jetstream.MessageHandler
	streamCfg   *jetstream.StreamConfig
	consumerCfg *jetstream.ConsumerConfig
}

func withSubj(subj string) SubscriptionOptions {
	return subscriptionCfgFn(func(opt *subscriptionOpt) error {
		opt.subj = subj
		return nil
	})
}

func withController(handler jetstream.MessageHandler) SubscriptionOptions {
	return subscriptionCfgFn(func(opt *subscriptionOpt) error {
		opt.msgHandler = handler
		return nil
	})
}

func withStreamCfg(streamCfg *jetstream.StreamConfig) SubscriptionOptions {
	return subscriptionCfgFn(func(opt *subscriptionOpt) error {
		opt.streamCfg = streamCfg

		if opt.streamCfg.MaxAge == 0 {
			opt.streamCfg.MaxAge = _msgMaxAge
		}
		return nil
	})
}

func withConsumerCfg(consumerCfg *jetstream.ConsumerConfig) SubscriptionOptions {
	return subscriptionCfgFn(func(opt *subscriptionOpt) error {
		opt.consumerCfg = consumerCfg

		if opt.consumerCfg.MaxDeliver == 0 {
			opt.consumerCfg.MaxDeliver = _maxDelivery
		}
		if opt.consumerCfg.AckWait == 0 {
			opt.consumerCfg.AckWait = _maxAckWait
		}
		return nil
	})
}

func WithAckWait(t time.Duration) SubscriptionOptions {
	return subscriptionCfgFn(func(opt *subscriptionOpt) error {
		opt.ackWait = t
		return nil
	})
}

func WithMaxDelivery(maxDelivery int) SubscriptionOptions {
	return subscriptionCfgFn(func(opt *subscriptionOpt) error {
		opt.maxDelivery = maxDelivery
		return nil
	})
}

func WithStreamMsgTTL(t time.Duration) SubscriptionOptions {
	return subscriptionCfgFn(func(opt *subscriptionOpt) error {
		opt.msgTTL = t
		return nil
	})
}
