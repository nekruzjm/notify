package telegram

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"

	"notifications/internal/service/telegram"
)

func (h *handler) Sent(msg jetstream.Msg) {
	h.logger.Info("msg Sent", zap.ByteString("data", msg.Data()))

	var (
		ctx     = context.Background()
		message struct {
			ChatID int64  `json:"chatID"`
			Text   string `json:"text"`
			Bot    string `json:"bot"`
		}
	)

	err := sonic.Unmarshal(msg.Data(), &message)
	if err != nil {
		h.logger.Error("sonic.Unmarshal error", zap.Error(err))
		return
	}

	err = msg.Ack()
	if err != nil {
		h.logger.Error("msg ack error", zap.Error(err))
		return
	}

	err = h.service.Send(ctx, telegram.Message{
		ChatID: message.ChatID,
		Text:   message.Text,
		Bot:    message.Bot,
	})
	if err != nil {
		h.logger.Error("Send error", zap.Error(err))
		return
	}
}
