package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/handler"
)

const AdditionalLocationMessage = " for save"

type Service interface {
	UpdateMessage(message *tgbotapi.Message) tgbotapi.MessageConfig
	UpdateLocation(message *tgbotapi.Message) tgbotapi.MessageConfig
	UpdateCallback(clb *tgbotapi.CallbackQuery) tgbotapi.MessageConfig
}

type Bot struct {
	handlers handler.Handlers
}

func New(handlers handler.Handlers) *Bot {
	return &Bot{handlers: handlers}
}

func (b Bot) UpdateMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	switch message.Text {
	case "🇺🇸 USA", "🇬🇧 UK", "🇨🇦 Canada", "🇫🇷 France", "🇩🇪 Germany", "🇯🇵 Japan":
		return b.handlers.HandleGetHolidays(message)
	case handler.StartMenu:
		return b.handlers.HandleStart(message)
	case handler.HolidayMenu:
		return b.handlers.HandleFlags(message)
	case handler.WeatherMenu:
		return b.handlers.HandleSendLocation(message, "")
	case handler.NotificationMenu:
		return b.handlers.HandleNotification(message)
	case handler.AddNotifyBtn:
		return b.handlers.HandleSendLocation(message, AdditionalLocationMessage)
	default:
		return tgbotapi.NewMessage(message.Chat.ID, "I don't know that command")
	}
}

func (b Bot) UpdateLocation(message *tgbotapi.Message) tgbotapi.MessageConfig {
	switch message.ReplyToMessage.Text {
	case handler.LocationMsg:
		return b.handlers.HandleGetWeatherByCoordinate(message)
	case handler.LocationMsg + AdditionalLocationMessage:
		return b.handlers.HandleCreateNotification(message)
	default:
		return tgbotapi.NewMessage(message.Chat.ID, "I don't now what to do with your location")
	}

}

func (b Bot) UpdateCallback(clb *tgbotapi.CallbackQuery) tgbotapi.MessageConfig {
	return b.handlers.HandleSaveTime(clb)
}
