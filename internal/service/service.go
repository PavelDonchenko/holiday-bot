package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handlers interface {
	HandleFlags(message *tgbotapi.Message) tgbotapi.MessageConfig
	HandleGetHolidays(message *tgbotapi.Message) tgbotapi.MessageConfig
}

type Bot struct {
	handlers Handlers
}

func New(handlers Handlers) *Bot {
	return &Bot{handlers: handlers}
}

func (b Bot) HandleMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	switch message.Text {
	case "🇺🇸 USA", "🇬🇧 UK", "🇨🇦 Canada", "🇫🇷 France", "🇩🇪 Germany", "🇯🇵 Japan":
		return b.handlers.HandleGetHolidays(message)
	default:
		return b.handlers.HandleFlags(message)
	}
}
