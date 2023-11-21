package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
)

type Service interface {
	HandleMessage(message *tgbotapi.Message) tgbotapi.MessageConfig
}

type Bot struct {
	handlers handler.Handlers
}

func New(handlers handler.Handlers) *Bot {
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
