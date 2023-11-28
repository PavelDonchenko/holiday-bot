package handler

import (
	"fmt"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
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
	sub := model.Subscription{
		ID:        uuid.New().String(),
		ChatID:    message.Chat.ID,
		Longitude: message.Location.Longitude,
		Latitude:  message.Location.Latitude,
	}

	err := h.db.Save(h.ctx, sub)
	if err != nil {
		h.log.Errorf("failed save location, err: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "error save location")
		return msg
	}

	return h.HandleShowTime(message)
}

func (h *Handler) HandleSaveTime(clb *tgbotapi.CallbackQuery) tgbotapi.MessageConfig {
	err := h.db.UpdateTime(h.ctx, clb.From.ID, clb.Data)
	if err != nil {
		h.log.Errorf("error save time, err: %v", err)
		return tgbotapi.NewMessage(clb.From.ID, "error save time")
	}

	return tgbotapi.NewMessage(clb.From.ID, "notification created successfully")
}
