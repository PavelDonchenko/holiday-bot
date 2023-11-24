package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
)

type Service interface {
	HandleMessage(message *tgbotapi.Message) tgbotapi.MessageConfig
	HandleForecastByLocation(message *tgbotapi.Message) tgbotapi.MessageConfig
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
	case handler.StartMenu:
		return b.handlers.HandleStart(message)
	case handler.HolidayMenu:
		return b.handlers.HandleFlags(message)
	case handler.WeatherMenu:
		return b.handlers.HandleWeatherCommand(message)
	default:
		return tgbotapi.NewMessage(message.Chat.ID, "I don't know that command")
	}
}
func (b Bot) HandleForecastByLocation(message *tgbotapi.Message) tgbotapi.MessageConfig {
	return b.handlers.HandleGetWeatherByCoordinate(message)
}
