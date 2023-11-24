package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
	"git.foxminded.ua/foxstudent106361/holiday-bot/internal/service"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	service service.Service
	cfg     config.Config
	log     *logrus.Logger
}

func New(api *tgbotapi.BotAPI, cfg config.Config, botService service.Service, log *logrus.Logger) *Bot {
	return &Bot{
		api:     api,
		cfg:     cfg,
		log:     log,
		service: botService,
	}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.cfg.Telegram.UpdateConfigTimeout

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		var msg tgbotapi.MessageConfig
		if update.Message != nil {
			if update.Message.Location != nil {
				msg = b.service.HandleForecastByLocation(update.Message)
				msg.ParseMode = tgbotapi.ModeHTML
			} else {
				msg = b.service.HandleMessage(update.Message)
			}

			_, err := b.api.Send(msg)
			if err != nil {
				b.log.Errorf("error sending message, err: %v", err)
				return
			}
		}
	}
}
