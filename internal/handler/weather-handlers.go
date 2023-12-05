package handler

import (
	"fmt"

	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) HandleSendLocation(message *tgbotapi.Message) tgbotapi.MessageConfig {
	btn := tgbotapi.KeyboardButton{
		Text:            "Send location",
		RequestLocation: true,
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, LocationMsg)

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})

	return msg
}

func (h *Handler) HandleGetWeatherByCoordinate(message *tgbotapi.Message) tgbotapi.MessageConfig {
	forecast, err := h.fetcher.GetForecast("", fmt.Sprint(message.Location.Longitude), fmt.Sprint(message.Location.Latitude))
	if err != nil {
		h.log.Error(err)
		return tgbotapi.MessageConfig{}
	}

	msg, err := utils.ParseForecast(*forecast)
	if err != nil {
		h.log.Error(err)
		return tgbotapi.MessageConfig{}
	}

	return tgbotapi.NewMessage(message.Chat.ID, msg)
}
