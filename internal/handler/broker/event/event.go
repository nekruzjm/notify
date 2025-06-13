package event

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"

	eventsrv "notifications/internal/service/event"
)

func (h *handler) Run(msg jetstream.Msg) {
	err := msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}

	h.service.RunJob()
}

func (h *handler) TopicSubscribed(msg jetstream.Msg) {
	var (
		ctx   = context.Background()
		event = new(eventsrv.Event)
	)

	err := sonic.Unmarshal(msg.Data(), event)
	if err != nil {
		h.logger.Error("msg unmarshal error", zap.Error(err))
		return
	}

	h.logger.Info("TopicSubscribed in process", zap.Any("eventID", event.ID))

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}

	if event.SubscribeAll {
		err = h.service.SubscribeAllUsers(ctx, event)
	} else {
		err = h.service.SubscribeUsers(ctx, event)
	}
	if err != nil {
		h.logger.Error("msg subscribe users error", zap.Error(err))
		return
	}
}

func (h *handler) TopicUnsubscribed(msg jetstream.Msg) {
	var (
		ctx   = context.Background()
		event = new(eventsrv.Event)
	)

	err := sonic.Unmarshal(msg.Data(), event)
	if err != nil {
		h.logger.Error("msg unmarshal error", zap.Error(err))
		return
	}

	h.logger.Info("TopicUnsubscribed in process", zap.Any("eventID", event.ID))

	err = h.service.UnsubscribeUsers(ctx, event)
	if err != nil {
		h.logger.Error("msg unsubscribe users error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}
}
