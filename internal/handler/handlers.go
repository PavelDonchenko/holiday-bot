package handler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/model"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/storage"

	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/client"
)

const (
	StartMenu       = "/start"
	HolidayMenu     = "/holiday"
	WeatherMenu     = "/weather"
	SubscribeMenu   = "/subscribe"
	UnsubscribeMenu = "/unsubscribe"
	AddNotifyBtn    = "Add notification"
	UpdateNotifyBtn = "Update notification"
	DeleteNotifyBtn = "Delete notification"

	LocationMsg = "Please send your location"
)

var _ Handlers = (*Handler)(nil)

type Handlers interface {
	HandleStart(message *tgbotapi.Message) tgbotapi.MessageConfig

	HandleFlags(message *tgbotapi.Message) tgbotapi.MessageConfig
	HandleGetHolidays(message *tgbotapi.Message) tgbotapi.MessageConfig

	HandleSendLocation(message *tgbotapi.Message) tgbotapi.MessageConfig
	HandleGetWeatherByCoordinate(message *tgbotapi.Message) tgbotapi.MessageConfig

	HandleGetTime(chatID int64) tgbotapi.MessageConfig
	HandleNotification(message *tgbotapi.Message) tgbotapi.MessageConfig
	HandleCreateNotification(message *tgbotapi.Message) (string, error)
	HandleSaveTime(time string, id string) error

	HandleSendSubscriptions(message *tgbotapi.Message) tgbotapi.MessageConfig
	HandleDeleteSub(clb *tgbotapi.CallbackQuery) error
	HandleDeleteLastSubscription() error

	HandleGetLastSubscription() (model.Subscription, error)
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
	msg := tgbotapi.NewMessage(message.Chat.ID, "Press menu button to see command list")

	return msg
}
