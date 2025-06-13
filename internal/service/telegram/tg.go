package telegram

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"notifications/pkg/lib/notifier/telegram"
)

func (s *service) Send(ctx context.Context, message Message) error {
	if message.ChatID == 0 {
		return errors.New("token or chatID cannot be empty")
	}

	err := s.tg.Send(ctx, telegram.Message{
		ChatID: message.ChatID,
		Text:   message.Text,
		Bot:    message.Bot,
	})
	if err != nil {
		s.sentry.CaptureException(err)
		s.logger.Error("err occurred during send message", zap.Error(err), zap.Any("message", message))
		return err
	}

	return nil
}
