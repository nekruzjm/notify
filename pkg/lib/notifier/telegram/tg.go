package telegram

import (
	"context"
	"errors"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (t *tg) Send(_ context.Context, message Message) error {
	if t.bots[message.Bot] == nil {
		t.logger.Error("tg: bot not found", zap.Any("message", message))
		return errors.New("tg: bot is not initialized")
	}

	msg := tgbotapi.NewMessage(message.ChatID, message.Text)
	msg.ParseMode = _defaultParseMode

	_, err := t.bots[message.Bot].Send(msg)
	if err != nil {
		t.logger.Error("tg: cannot send message", zap.Error(err))
		return err
	}

	return nil
}
