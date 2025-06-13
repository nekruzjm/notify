package email

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"

	"notifications/internal/service/email"
)

func (h *handler) Sent(msg jetstream.Msg) {
	h.logger.Info("msg Sent", zap.ByteString("data", msg.Data()))

	var (
		ctx  = context.Background()
		data = struct {
			Body map[string]string `json:"body"`
		}{}
	)

	err := sonic.Unmarshal(msg.Data(), &data)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err))
		return
	}

	err = h.service.Send(ctx, email.Email{
		Body: data.Body,
	})
	if err != nil {
		h.logger.Error("Send error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}
}
