package handler

import (
	"fmt"
	"regexp"
	"strconv"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/utils"
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

func (h *Handler) HandleShowTime(chatID int64, id string) tgbotapi.MessageConfig {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for hour := 0; hour < 24; hour++ {
		buttonText := fmt.Sprintf("%02d:00", hour)
		btn := tgbotapi.NewInlineKeyboardButtonData(buttonText, fmt.Sprintf("%s&%s", buttonText, id))
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{btn})
	}

	replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboard...)

	msg := tgbotapi.NewMessage(chatID, "Please choose a notification time:")
	msg.ReplyMarkup = replyKeyboard

	return msg
}

func (h *Handler) HandleCreateNotification(message *tgbotapi.Message) (string, error) {
	sub := model.Subscription{
		ID:        uuid.New().String(),
		ChatID:    message.Chat.ID,
		Longitude: utils.Round(message.Location.Longitude, 2),
		Latitude:  utils.Round(message.Location.Latitude, 2),
	}

	id, err := h.db.Save(h.ctx, sub)
	if err != nil {
		h.log.Errorf("failed save location, err: %v", err)
		return "", err
	}

	return id, nil
}

func (h *Handler) HandleSaveTime(clb *tgbotapi.CallbackQuery) error {
	err := h.db.UpdateTime(h.ctx, clb.Data)
	if err != nil {
		h.log.Errorf("error save time, err: %v", err)
		return err
	}

	return nil
}

func (h *Handler) HandleSendSubscriptions(message *tgbotapi.Message) tgbotapi.MessageConfig {
	subs, _ := h.db.GetSubscriptions(h.ctx, message.Chat.ID)

	var keyboard [][]tgbotapi.InlineKeyboardButton
	var msg tgbotapi.MessageConfig
	if subs != nil {
		for _, sub := range subs {
			buttonText := fmt.Sprintf("Longitude:%v, Latitude:%v, time:%s", sub.Longitude, sub.Latitude, sub.NotifyTime)
			btn := tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonText)
			keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{btn})
		}

		replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboard...)

		msg = tgbotapi.NewMessage(message.Chat.ID, "Please choose a notification:")
		msg.ReplyMarkup = replyKeyboard
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "You don't have any subscription")
	}

	return msg
}

func (h *Handler) HandleDeleteSub(clb *tgbotapi.CallbackQuery) error {
	long, lat, time, err := parseLocationTime(clb.Data)
	if err != nil {
		return err
	}

	err = h.db.DeleteSubscription(h.ctx, long, lat, time)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleGetSubscriptionID(clb *tgbotapi.CallbackQuery) (string, error) {
	long, lat, time, err := parseLocationTime(clb.Data)
	if err != nil {
		return "", err
	}

	id, err := h.db.GetSubscriptionByID(h.ctx, long, lat, time)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (h *Handler) HandleUpdateTime(time string) error {

	err := h.db.UpdateTime(h.ctx, time)
	if err != nil {
		return err
	}
	return nil
}

func parseLocationTime(input string) (float64, float64, string, error) {
	pattern := `Longitude:\s*([-+]?\d*\.\d+|\d+),\s*Latitude:\s*([-+]?\d*\.\d+|\d+),\s*time:\s*([0-9]{2}:[0-9]{2})`

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(input)

	if len(matches) != 4 {
		return 0, 0, "", fmt.Errorf("invalid input format")
	}

	longitude, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, 0, "", err
	}

	latitude, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return 0, 0, "", err
	}

	timeValue := matches[3]

	return longitude, latitude, timeValue, nil
}
