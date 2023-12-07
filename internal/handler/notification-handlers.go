package handler

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"gopkg.in/ugjka/go-tz.v2/tz"
)

func (h *Handler) HandleGetTime(chatID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "Please type the time you want to receive notification (IN '13:00' FORMAT):")

	return msg
}

func (h *Handler) HandleDeleteLastSubscription() error {
	return h.db.DeleteLastSubscription(h.ctx)
}

func (h *Handler) HandleCreateNotification(message *tgbotapi.Message) (string, error) {
	sub := model.Subscription{
		ID:        uuid.New().String(),
		ChatID:    message.Chat.ID,
		Longitude: utils.Round(message.Location.Longitude, 2),
		Latitude:  utils.Round(message.Location.Latitude, 2),
		CreatedAt: time.Now().UTC(),
	}

	id, err := h.db.Save(h.ctx, sub)
	if err != nil {
		h.log.Errorf("failed save location, err: %v", err)
		return "", err
	}

	return id, nil
}

func (h *Handler) HandleSaveTime(timeToSave string, sub model.Subscription) error {
	t, err := timeToUTC(timeToSave, sub.Latitude, sub.Longitude)
	if err != nil {
		h.log.Errorf("error parse time, err: %v", err)
		return err
	}
	if err := h.db.UpdateTime(h.ctx, t, sub.ID); err != nil {
		h.log.Errorf("error save time, err: %v", err)
		return err
	}

	return nil
}

func (h *Handler) HandleSendSubscriptions(message *tgbotapi.Message) tgbotapi.MessageConfig {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var msg tgbotapi.MessageConfig

	subs, err := h.db.GetSubscriptions(h.ctx, message.Chat.ID)
	if err != nil {
		msg = tgbotapi.NewMessage(message.Chat.ID, "failed get subscriptions")
		return msg
	}

	if len(subs) > 0 {
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
	long, lat, notificationTime, err := parseLocationTime(clb.Data)
	if err != nil {
		return err
	}

	return h.db.DeleteSubscription(h.ctx, long, lat, notificationTime)
}

func (h *Handler) HandleGetLastSubscription() (model.Subscription, error) {
	sub, err := h.db.GetLastSubscription(h.ctx)
	if err != nil {
		return model.Subscription{}, err
	}

	return sub, nil
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

func timeToUTC(userTime string, lat, lon float64) (time.Time, error) {
	zone, err := tz.GetZone(tz.Point{
		Lon: lon, Lat: lat,
	})
	if err != nil {
		return time.Time{}, err
	}

	userTimeParsed, err := time.Parse("15:04", userTime)
	if err != nil {
		return time.Time{}, err
	}

	location, err := time.LoadLocation(zone[0])
	if err != nil {
		return time.Time{}, err
	}

	_, offset := time.Now().In(location).Zone()
	loc := time.FixedZone(zone[0], offset)

	userTimeWithZone := time.Date(0, 1, 1, userTimeParsed.Hour(), userTimeParsed.Minute(), 0, 0, loc)

	utcTime := userTimeWithZone.UTC()

	return utcTime, nil
}
