package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"notifications/pkg/lib/config"
	"notifications/pkg/lib/observer/logger"
)

var Module = fx.Provide(New)

type Telegram interface {
	Send(context.Context, Message) error
}

type Params struct {
	fx.In

	Logger logger.Logger
	Config config.Config
}

type tg struct {
	logger logger.Logger
	bots   map[string]*tgbotapi.BotAPI
}

func New(p Params) Telegram {
	return &tg{
		logger: p.Logger,
		bots:   _connect(p),
	}
}

const (
	_tcbTransferBot        = "tcbTransferBot"
	_dbStatBot             = "dbStatBot"
	_providerSuggestionBot = "providerSuggestionBot"
)

func _connect(p Params) map[string]*tgbotapi.BotAPI {
	var (
		keys = []string{_tcbTransferBot, _dbStatBot, _providerSuggestionBot}
		bots = make(map[string]*tgbotapi.BotAPI, len(keys))
	)

	for _, key := range keys {
		token := p.Config.GetString("telegram." + key)
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			p.Logger.Error("cannot connect to bot", zap.Error(err), zap.String("bot", key))
			continue
		}
		bots[key] = bot
	}

	p.Logger.Info("Successfully connected to bots")

	return bots
}
