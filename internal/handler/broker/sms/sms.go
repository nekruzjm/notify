package sms

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (h *handler) Sent(msg jetstream.Msg) {
	h.logger.Info("msg Sent", zap.ByteString("data", msg.Data()))

	var (
		ctx  = context.Background()
		body = make(map[string]string)
	)

	err := sonic.Unmarshal(msg.Data(), &body)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}

	err = h.service.Send(ctx, body)
	if err != nil {
		h.logger.Error("Send error", zap.Error(err))
		return
	}
}
