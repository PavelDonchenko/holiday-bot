package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) HandleNotification(message *tgbotapi.Message) tgbotapi.MessageConfig {
	menuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(AddNotifyBtn),
			tgbotapi.NewKeyboardButton(UpdateNotifyBtn),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DeleteNotifyBtn),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, "What do you want:")
	msg.ReplyMarkup = menuKeyboard

	return msg
}

func (h *Handler) HandleShowTime(message *tgbotapi.Message) tgbotapi.MessageConfig {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for hour := 0; hour < 24; hour++ {
		buttonText := fmt.Sprintf("%02d:00", hour)
		btn := tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonText)
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{btn})
	}

	replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Please choose a notification time:")
	msg.ReplyMarkup = replyKeyboard

	return msg
}

func (h *Handler) HandleCreateNotification(message *tgbotapi.Message) tgbotapi.MessageConfig {
	panic("dasd")
}
