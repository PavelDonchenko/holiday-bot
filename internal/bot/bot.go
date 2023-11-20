package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"git.foxminded.ua/foxstudent106361/holiday-bot/pkg/logging"

	"git.foxminded.ua/foxstudent106361/holiday-bot/config"
)

type Service interface {
	HandleMessage(message *tgbotapi.Message) tgbotapi.MessageConfig
}

type Bot struct {
	api     *tgbotapi.BotAPI
	service Service
	cfg     config.Config
	log     logging.Logger
}

func New(api *tgbotapi.BotAPI, cfg config.Config, service Service, log logging.Logger) *Bot {
	return &Bot{
		api:     api,
		cfg:     cfg,
		log:     log,
		service: service,
	}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.cfg.Telegram.UpdateConfigTimeout

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			msg := b.service.HandleMessage(update.Message)

			_, err := b.api.Send(msg)
			if err != nil {
				b.log.Errorf("error sending message, err: %v", err)
				return
			}
		}
	}
}
