package handler

import (
	"context"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/client"
)

const (
	StartMenu        = "/start"
	HolidayMenu      = "🏝 Holiday"
	WeatherMenu      = "☀️ Weather"
	NotificationMenu = "⏰ Notification"
	AddNotifyBtn     = "Add notification"
	UpdateNotifyBtn  = "Update notification"
	DeleteNotifyBtn  = "Delete notification"

	LocationMsg = "Please send your location"
)

var _ Handlers = (*Handler)(nil)

type Handlers interface {
	HandleStart(message *tgbotapi.Message) tgbotapi.MessageConfig

	HandleFlags(message *tgbotapi.Message) tgbotapi.MessageConfig
	HandleGetHolidays(message *tgbotapi.Message) tgbotapi.MessageConfig

	HandleSendLocation(message *tgbotapi.Message, addMsg string) tgbotapi.MessageConfig
	HandleGetWeatherByCoordinate(message *tgbotapi.Message) tgbotapi.MessageConfig

	HandleShowTime(update *tgbotapi.Message) tgbotapi.MessageConfig
	HandleNotification(message *tgbotapi.Message) tgbotapi.MessageConfig
	HandleCreateNotification(message *tgbotapi.Message) tgbotapi.MessageConfig
	HandleSaveTime(clb *tgbotapi.CallbackQuery) tgbotapi.MessageConfig
}

type Handler struct {
	log     *logrus.Logger
	fetcher client.Fetcher
	db      storage.Storage
	ctx     context.Context
}

func New(ctx context.Context, log *logrus.Logger, fetcher client.Fetcher, db storage.Storage) *Handler {
	return &Handler{
		log:     log,
		fetcher: fetcher,
		db:      db,
		ctx:     ctx,
	}
}

func (h *Handler) HandleStart(message *tgbotapi.Message) tgbotapi.MessageConfig {
	menuKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(WeatherMenu),
			tgbotapi.NewKeyboardButton(HolidayMenu),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(NotificationMenu),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Here is a list of what I can do:")
	msg.ReplyMarkup = menuKeyboard

	return msg
}
